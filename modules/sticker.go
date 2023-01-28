package modules

import (
	"bytes"
	"context"
	"gowa/utils"
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

func sticker_callback(cli *whatsmeow.Client, msg *events.Message) {
	if msg.Message.ExtendedTextMessage == nil || msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage.ImageMessage == nil {
		cli.SendMessage(context.Background(), msg.Info.Chat,
			utils.NewMessage("Reply to a message with photo dipshit!", msg),
		)
		return
	}

	downloaded_image, err := cli.Download(msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage.ImageMessage)

	if err != nil {
		cli.SendMessage(context.Background(), msg.Info.Chat,
			utils.NewMessage("Failed to download photo!", msg),
		)
		return
	}
	decoder := getDecoder(msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage.ImageMessage.Mimetype)
	if decoder == nil {
		cli.SendMessage(context.Background(), msg.Info.Chat,
			utils.NewMessage("Invalid MimeType!", msg),
		)
		return
	}
	img, err := decoder(bytes.NewReader(downloaded_image))
	if err != nil {
		cli.SendMessage(context.Background(), msg.Info.Chat,
			utils.NewMessage("Failed to decode photo!", msg),
		)
		return
	}
	webpByte, err := webp.EncodeRGBA(img, *proto.Float32(1))
	if err != nil {
		cli.SendMessage(context.Background(), msg.Info.Chat,
			utils.NewMessage("Failed to encode sticker!", msg),
		)
		return
	}

	uploadImage, err := cli.Upload(context.Background(), webpByte, whatsmeow.MediaImage)
	if err != nil {
		cli.SendMessage(context.Background(), msg.Info.Chat,
			utils.NewMessage("Failed to upload sticker!", msg),
		)
		return
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
}

var Sticker Command = Command{
	".sticker",
	sticker_callback,
	sticker_help,
}
