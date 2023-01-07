package api

import (
	"context"
	"fmt"

	"github.com/ncyellow/GophKeeper/internal/client/config"
	"github.com/ncyellow/GophKeeper/internal/models"
	"github.com/ncyellow/GophKeeper/internal/server/gprcserver/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// GRPCSender структура http клиента
type GRPCSender struct {
	conn   *grpc.ClientConn
	client proto.GophKeeperServerClient
	conf   *config.Config
	userID *int64
}

func (g *GRPCSender) Register(login string, pwd string) error {
	response, err := g.client.Register(context.Background(), &proto.RegisterRequest{
		Login:    login,
		Password: pwd,
	})

	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.AlreadyExists {
				return ErrUserAlreadyExists
			}
		}
		return ErrInternalServer
	}

	userID := response.GetUser()
	g.userID = &userID

	fmt.Printf("userid = %d", userID)
	fmt.Printf("userid = %d", *g.userID)
	return nil
}

func (g *GRPCSender) SignIn(login string, pwd string) error {
	response, err := g.client.SignIn(context.Background(), &proto.RegisterRequest{
		Login:    login,
		Password: pwd,
	})

	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.AlreadyExists {
				return ErrUserAlreadyExists
			}
		}
		return ErrInternalServer
	}

	userID := response.GetUser()
	g.userID = &userID
	fmt.Printf("userid = %d", userID)
	fmt.Printf("userid = %d", *g.userID)
	return nil
}

func (g *GRPCSender) AddCard(card *models.Card) error {

	if g.userID == nil {
		return ErrAuthRequire
	}

	_, err := g.client.AddCard(context.Background(), &proto.AddCardRequest{
		Card: &proto.Card{
			Id:       card.ID,
			Fio:      card.FIO,
			Number:   card.Number,
			Date:     card.Date,
			Cvv:      card.CVV,
			Metainfo: card.MetaInfo,
		},
		User: *g.userID,
	})

	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.AlreadyExists {
				return ErrAlreadyExists
			}
		}
		return ErrInternalServer
	}
	return nil
}

func (g *GRPCSender) Card(cardID string) (*models.Card, error) {
	if g.userID == nil {
		return nil, ErrAuthRequire
	}

	response, err := g.client.Card(context.Background(), &proto.CardRequest{
		Id:   cardID,
		User: *g.userID,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.NotFound {
				return nil, ErrNotFound
			}
		}
		return nil, ErrInternalServer
	}

	respCard := response.GetCard()
	return &models.Card{
		ID:       respCard.GetId(),
		FIO:      respCard.GetFio(),
		Number:   respCard.GetNumber(),
		Date:     respCard.GetDate(),
		CVV:      respCard.GetCvv(),
		MetaInfo: respCard.GetMetainfo(),
	}, nil
}

func (g *GRPCSender) DelCard(cardID string) error {
	if g.userID == nil {
		return ErrAuthRequire
	}
	_, err := g.client.DeleteCard(context.Background(), &proto.DeleteCardRequest{
		Id:   cardID,
		User: *g.userID,
	})
	if err != nil {
		return ErrInternalServer
	}
	return nil
}

func (g *GRPCSender) AddLogin(login *models.Login) error {
	if g.userID == nil {
		return ErrAuthRequire
	}

	_, err := g.client.AddLogin(context.Background(), &proto.AddLoginRequest{
		Login: &proto.Login{
			Id:       login.ID,
			Login:    login.Login,
			Password: login.Password,
			Metainfo: login.MetaInfo,
		},
		User: *g.userID,
	})

	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.AlreadyExists {
				return ErrAlreadyExists
			}
		}
		return ErrInternalServer
	}
	return nil
}

func (g *GRPCSender) Login(loginID string) (*models.Login, error) {
	if g.userID == nil {
		return nil, ErrAuthRequire
	}

	response, err := g.client.Login(context.Background(), &proto.LoginRequest{
		Id:   loginID,
		User: *g.userID,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.NotFound {
				return nil, ErrNotFound
			}
		}
		return nil, ErrInternalServer
	}

	respLogin := response.GetLogin()
	return &models.Login{
		ID:       respLogin.GetId(),
		Login:    respLogin.GetLogin(),
		Password: respLogin.GetPassword(),
		MetaInfo: respLogin.GetMetainfo(),
	}, nil
}

func (g *GRPCSender) DelLogin(loginID string) error {
	if g.userID == nil {
		return ErrAuthRequire
	}
	_, err := g.client.DeleteLogin(context.Background(), &proto.DeleteLoginRequest{
		Id:   loginID,
		User: *g.userID,
	})
	if err != nil {
		return ErrInternalServer
	}
	return nil
}

func (g *GRPCSender) AddText(text *models.Text) error {
	if g.userID == nil {
		return ErrAuthRequire
	}

	_, err := g.client.AddText(context.Background(), &proto.AddTextRequest{
		Text: &proto.Text{
			Id:       text.ID,
			Content:  text.Content,
			Metainfo: text.MetaInfo,
		},
		User: *g.userID,
	})

	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.AlreadyExists {
				return ErrAlreadyExists
			}
		}
		return ErrInternalServer
	}
	return nil
}

func (g *GRPCSender) Text(textID string) (*models.Text, error) {
	if g.userID == nil {
		return nil, ErrAuthRequire
	}

	response, err := g.client.Text(context.Background(), &proto.TextRequest{
		Id:   textID,
		User: *g.userID,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.NotFound {
				return nil, ErrNotFound
			}
		}
		return nil, ErrInternalServer
	}

	respText := response.GetText()
	return &models.Text{
		ID:       respText.GetId(),
		Content:  respText.GetContent(),
		MetaInfo: respText.GetMetainfo(),
	}, nil
}

func (g *GRPCSender) DelText(textID string) error {
	if g.userID == nil {
		return ErrAuthRequire
	}
	_, err := g.client.DeleteText(context.Background(), &proto.DeleteTextRequest{
		Id:   textID,
		User: *g.userID,
	})
	if err != nil {
		return ErrInternalServer
	}
	return nil
}

func (g *GRPCSender) AddBin(binary *models.Binary) error {
	if g.userID == nil {
		return ErrAuthRequire
	}

	_, err := g.client.AddBinary(context.Background(), &proto.AddBinRequest{
		Binary: &proto.Binary{
			Id:       binary.ID,
			Data:     binary.Data,
			Metainfo: binary.MetaInfo,
		},
		User: *g.userID,
	})

	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.AlreadyExists {
				return ErrAlreadyExists
			}
		}
		return ErrInternalServer
	}
	return nil
}

func (g *GRPCSender) Bin(binID string) (*models.Binary, error) {
	if g.userID == nil {
		return nil, ErrAuthRequire
	}

	response, err := g.client.Binary(context.Background(), &proto.BinRequest{
		Id:   binID,
		User: *g.userID,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.NotFound {
				return nil, ErrNotFound
			}
		}
		return nil, ErrInternalServer
	}

	respText := response.GetBinary()
	return &models.Binary{
		ID:       respText.GetId(),
		Data:     respText.GetData(),
		MetaInfo: respText.GetMetainfo(),
	}, nil
}

func (g *GRPCSender) DelBin(binID string) error {
	if g.userID == nil {
		return ErrAuthRequire
	}
	_, err := g.client.DeleteBinary(context.Background(), &proto.DeleteBinRequest{
		Id:   binID,
		User: *g.userID,
	})
	if err != nil {
		return ErrInternalServer
	}
	return nil
}

func NewGRPCSender(conf *config.Config) *GRPCSender {
	// устанавливаем соединение с сервером
	conn, err := grpc.Dial(conf.GRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Err(err)
	}
	client := proto.NewGophKeeperServerClient(conn)
	return &GRPCSender{
		conf:   conf,
		conn:   conn,
		client: client,
	}
}
