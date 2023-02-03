package modules

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"

	"github.com/chai2010/webp"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

const sticker_help string = "Convert any given image into a sticker."

func getDecoder(mimeType *string) func(r io.Reader) (image.Image, error) {
	if *mimeType == "image/jpeg" {
		return jpeg.Decode
	}
	if *mimeType == "image/png" {
		return png.Decode
	}
	if *mimeType == "image/webp" {
		return webp.Decode
	}

	return nil
}

func sticker_callback(cli *whatsmeow.Client, msg *events.Message) error {
	if msg.Message.ExtendedTextMessage == nil || msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage.ImageMessage == nil {
		return errors.New("Reply to a message with photo! QuotedMessage / Image was nil.")
	}

	downloaded_image, err := cli.Download(msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage.ImageMessage)

	if err != nil {
		return errors.New("Failed to download photo!")
	}
	decoder := getDecoder(msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage.ImageMessage.Mimetype)
	if decoder == nil {
		return errors.New("Invalid MimeType")
	}
	img, err := decoder(bytes.NewReader(downloaded_image))
	if err != nil {
		return errors.New("Failed to decode photo!")
	}
	webpByte, err := webp.EncodeRGBA(img, *proto.Float32(1))
	if err != nil {
		return errors.New("Failed to encode sticker!")
	}

	uploadImage, err := cli.Upload(context.Background(), webpByte, whatsmeow.MediaImage)
	if err != nil {
		return errors.New("Failed to upload sticker!")
	}

	// Showed frame with thumbnail
	result := &waProto.Message{
		StickerMessage: &waProto.StickerMessage{
			Url:               proto.String(uploadImage.URL),
			FileSha256:        uploadImage.FileSHA256,
			FileEncSha256:     uploadImage.FileEncSHA256,
			MediaKey:          uploadImage.MediaKey,
			Mimetype:          proto.String(http.DetectContentType(webpByte)),
			DirectPath:        proto.String(uploadImage.DirectPath),
			FileLength:        proto.Uint64(uint64(len(webpByte))),
			FirstFrameSidecar: webpByte,
			PngThumbnail:      webpByte,
			ContextInfo:       &waProto.ContextInfo{StanzaId: &msg.Info.ID, Participant: proto.String(msg.Info.MessageSource.Sender.String()), QuotedMessage: msg.Message},
		},
	}

	cli.SendMessage(context.Background(), msg.Info.Chat, result)

	return nil
}

var Sticker Command = Command{
	".sticker",
	sticker_callback,
	sticker_help,
}
