package generate_image

import (
	"fmt"
	"image"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/paavill/awesome-tagger-bot/domain/cases/send_images"
	"github.com/paavill/awesome-tagger-bot/domain/context"
	"github.com/paavill/awesome-tagger-bot/domain/state_machine"
)

func New() state_machine.State {
	return &initState{}
}

type initState struct {
	state_machine.Dumper
}

func (s *initState) ProcessCallbackRequest(ctx context.Context, callback *tgbotapi.CallbackQuery) (state_machine.ProcessResponse, error) {
	return nil, nil
}

func (s *initState) ProcessMessage(ctx context.Context, message *tgbotapi.Message) (state_machine.ProcessResponse, error) {
	return Run(ctx, message)
}

func Run(ctx context.Context, message *tgbotapi.Message) (state_machine.ProcessResponse, error) {
	if message == nil {
		return nil, fmt.Errorf("[generate_image] message is nil")
	}

	selfName := ctx.Services().Bot().Self.UserName
	if message.Text != "/generate_image" && message.Text != "/generate_image@"+selfName {
		return nil, nil
	}

	ctx.Logger().Info("[generate_image] start")
	defer ctx.Logger().Info("[generate_image] finish")

	messageChat := message.Chat
	if messageChat == nil {
		return nil, fmt.Errorf("[generate_image] message chat is nil")
	}

	newMessage := tgbotapi.NewMessage(messageChat.ID, "Что нарисовать?")
	_, err := ctx.Services().Bot().Send(newMessage)
	if err != nil {
		return nil, fmt.Errorf("[generate_image] failed to send message: %s", err)
	}

	return state_machine.NewProcessResponse(false, &generateImageState{}), nil
}

type generateImageState struct {
	state_machine.Dumper
}

func (s *generateImageState) ProcessCallbackRequest(ctx context.Context, callback *tgbotapi.CallbackQuery) (state_machine.ProcessResponse, error) {
	return nil, nil
}

func (s *generateImageState) ProcessMessage(ctx context.Context, message *tgbotapi.Message) (state_machine.ProcessResponse, error) {
	return runGeneration(ctx, message)
}

func runGeneration(ctx context.Context, message *tgbotapi.Message) (state_machine.ProcessResponse, error) {
	ctx.Logger().Info("[generate_image] start generation")
	defer ctx.Logger().Info("[generate_image] finish generation")

	if message == nil {
		return nil, fmt.Errorf("[generate_image] message is nil")
	}

	query := message.Text
	if query == "" {
		return nil, fmt.Errorf("[generate_image] query is empty")
	}

	img, err := ctx.Services().Kandinsky().GenerateImage(query)
	if err != nil {
		return nil, fmt.Errorf("[generate_image] failed to generate image: %s", err)
	}

	messageChat := message.Chat
	if messageChat == nil {
		return nil, fmt.Errorf("[generate_image] message chat is nil")
	}

	err = send_images.Run(ctx, messageChat.ID, "", []*image.Image{img})
	if err != nil {
		return nil, fmt.Errorf("[generate_image] failed to send image: %s", err)
	}
	return nil, nil
}
