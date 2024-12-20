package modules

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gowa/utils"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

// var INVIDIOUS_HOST string = "https://inv.riverside.rocks/"
var INVIDIOUS_HOST string = "https://y.com.sb/"
var INVIDIOUS_VIDEO_PATH string = "api/v1/videos/"

type AdaptiveFormat struct {
	Bitrate         string `json:"bitrate"`
	Url             string `json:"url"`
	Type            string `json:"type"`
	Encoding        string `json:"encoding"`
	AudioSamplerate int    `json:"audioSamplerate"`
	Container       string `json:"container"`
	AudioChannels   int    `json:"audioChannels"`
}

type FormatStream struct {
	Url          string `json:"url"`
	Itag         string `json:"itag"`
	Type         string `json:"type"`
	Quality      string `json:"quality"`
	Fps          int    `json:"fps"`
	Container    string `json:"container"`
	Encoding     string `json:"encoding"`
	Resolution   string `json:"resolution"`
	QualityLabel string `json:"qualityLabel"`
	Size         string `json:"size"`
}

type VideoDataResponse struct {
	Title           string           `json:"title"`
	AdaptiveFormats []AdaptiveFormat `json:"adaptiveFormats"`
	FormatStreams   []FormatStream   `json:"formatStreams"`
}

var ytdl_help string = "for downloading youtube videos/audio"

// HELPER FUNCTIONS
func last[T any](arr []T) T {
	return arr[len(arr)-1]
}

func filter[T comparable](arr []T, filterFunc func(el T) bool) {
	n := 0
	for _, val := range arr {
		if filterFunc(val) {
			arr[n] = val
			n++
		}
	}
	arr = arr[:n]
}

func ytdl_callback(cli *whatsmeow.Client, msg *events.Message) error {

	text := utils.GetText(msg)
	splitFields := strings.Fields(*text)
	if len(splitFields) < 2 {
		return errors.New("enter a url")
	}

	if !(utils.Contains(splitFields, "-a") || utils.Contains(splitFields, "-v")) {
		return errors.New("Didnt match any path, use flags like -a -v")
	}

	u, err := url.Parse(splitFields[1])
	if err != nil {
		return errors.New("error parsing url")
	}

	resp, err := http.Get(INVIDIOUS_HOST + INVIDIOUS_VIDEO_PATH + u.Query()["v"][0])

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	var jsonResp VideoDataResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Error reading resp body")
	}

	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		return errors.New("Deserializing failed")
	}
	// Deserializing DONE
	if utils.Contains(splitFields, "-a") {
		audioDownload(cli, msg, &jsonResp)
		return nil
	}

	if utils.Contains(splitFields, "-v") {
		return errors.New("Not Implemented yet!")
	}

	return errors.New("Didnt match any path, use flags like -a -v")
}

func audioDownload(cli *whatsmeow.Client, msg *events.Message, jsonResp *VideoDataResponse) error {
	filter(jsonResp.AdaptiveFormats, func(i AdaptiveFormat) bool {
		return strings.HasPrefix(i.Type, "audio")
	})

	sort.Slice(jsonResp.AdaptiveFormats[:], func(i, j int) bool {
		return jsonResp.AdaptiveFormats[i].AudioSamplerate < jsonResp.AdaptiveFormats[j].AudioSamplerate
	})

	last_el := last(jsonResp.AdaptiveFormats)
	fmt.Println(last_el)

	fmt.Println("downloading audio file")
	client := &http.Client{}
	req, _ := http.NewRequest("GET", last_el.Url, nil)
	req.Header.Set("range", "bytes=0-") // else google limits dl speed to 6kbps
	dataResp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer dataResp.Body.Close()
	fmt.Println("Reading all bytes")
	dataBytes, err := io.ReadAll(dataResp.Body)
	fmt.Println("Reading done")
	if err != nil {
		return err
	}

	fmt.Println("Uploading audio file")
	audioUploaded, err := cli.Upload(context.Background(), dataBytes, whatsmeow.MediaDocument)

	if err != nil {
		return err
	}

	/*
	* ONLY 2 things work
	* audio/mpeg -> mp3
	* audio/ogg; codecs=opus -> ogg (voice message)
	* LITERALLY ANYTHING ELSE does not seem to work on android client
	* need some lib/ffmpeg exec (last resort) to convert encoding/container
	*
	* Currently sending the audio file as document, too lazy to convert :P
	 */

	// Compose WhatsApp Proto
	content := &waE2E.Message{
		DocumentMessage: &waE2E.DocumentMessage{
			FileName:   proto.String(jsonResp.Title + "." + last_el.Container),
			URL:        proto.String(audioUploaded.URL),
			DirectPath: proto.String(audioUploaded.DirectPath),
			// Mimetype:   proto.String(strings.ReplaceAll(last_el.Type, "\"", "")),
			Mimetype:      proto.String(http.DetectContentType(dataBytes)),
			FileLength:    proto.Uint64(audioUploaded.FileLength),
			FileSHA256:    audioUploaded.FileSHA256,
			FileEncSHA256: audioUploaded.FileEncSHA256,
			MediaKey:      audioUploaded.MediaKey,
			ContextInfo:   &waE2E.ContextInfo{StanzaID: &msg.Info.ID, Participant: proto.String(msg.Info.MessageSource.Sender.String()), QuotedMessage: msg.Message},
		},
	}

	_, err = cli.SendMessage(context.Background(), msg.Info.Chat, content)

	return err
}

var Ytdl Command = Command{
	".ytdl",
	ytdl_callback,
	ytdl_help,
}
