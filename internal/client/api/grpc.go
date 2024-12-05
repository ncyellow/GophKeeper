package api

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/ncyellow/GophKeeper/internal/client/config"
	"github.com/ncyellow/GophKeeper/internal/models"
	proto2 "github.com/ncyellow/GophKeeper/internal/proto"
)

// GRPCSender structure of the grpc client. Implements the Sender interface. See the respective methods for all comments.
type GRPCSender struct {
	conn   *grpc.ClientConn
	client proto2.GophKeeperServerClient
	conf   *config.Config
	userID *int64
}

// NewGRPCSender constructor
func NewGRPCSender(conf *config.Config) (*GRPCSender, error) {
	// establish a connection to the server
	conn, err := grpc.Dial(conf.GRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := proto2.NewGophKeeperServerClient(conn)
	return &GRPCSender{
		conf:   conf,
		conn:   conn,
		client: client,
	}, nil
}

func (g *GRPCSender) Register(login string, pwd string) error {
	response, err := g.client.Register(context.Background(), &proto2.RegisterRequest{
		Login:    login,
		Password: pwd,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.AlreadyExists {
				return fmt.Errorf(FmtErrUserAlreadyExists, err)
			}
		}
		return fmt.Errorf(FmtErrInternalServer, err)
	}
	userID := response.GetUser()
	g.userID = &userID

	fmt.Printf("userid = %d", userID)
	fmt.Printf("userid = %d", *g.userID)
	return nil
}

func (g *GRPCSender) SignIn(login string, pwd string) error {
	response, err := g.client.SignIn(context.Background(), &proto2.RegisterRequest{
		Login:    login,
		Password: pwd,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.AlreadyExists {
				return fmt.Errorf(FmtErrUserAlreadyExists, err)
			}
		}
		return fmt.Errorf(FmtErrInternalServer, err)
	}
	userID := response.GetUser()
	g.userID = &userID
	return nil
}

func (g *GRPCSender) AddCard(card *models.Card) error {
	if g.userID == nil {
		return ErrAuthRequire
	}
	_, err := g.client.AddCard(context.Background(), &proto2.AddCardRequest{
		Card: &proto2.Card{
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
				return fmt.Errorf(FmtErrAlreadyExists, err)
			}
		}
		return fmt.Errorf(FmtErrInternalServer, err)
	}
	return nil
}

func (g *GRPCSender) Card(cardID string) (*models.Card, error) {
	if g.userID == nil {
		return nil, ErrAuthRequire
	}

	response, err := g.client.Card(context.Background(), &proto2.CardRequest{
		Id:   cardID,
		User: *g.userID,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.NotFound {
				return nil, fmt.Errorf(FmtErrNotFound, err)
			}
		}
		return nil, fmt.Errorf(FmtErrInternalServer, err)
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
	_, err := g.client.DeleteCard(context.Background(), &proto2.DeleteCardRequest{
		Id:   cardID,
		User: *g.userID,
	})
	if err != nil {
		return fmt.Errorf(FmtErrInternalServer, err)
	}
	return nil
}

func (g *GRPCSender) AddLogin(login *models.Login) error {
	if g.userID == nil {
		return ErrAuthRequire
	}

	_, err := g.client.AddLogin(context.Background(), &proto2.AddLoginRequest{
		Login: &proto2.Login{
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
				return fmt.Errorf(FmtErrAlreadyExists, err)
			}
		}
		return fmt.Errorf(FmtErrInternalServer, err)
	}
	return nil
}

func (g *GRPCSender) Login(loginID string) (*models.Login, error) {
	if g.userID == nil {
		return nil, ErrAuthRequire
	}

	response, err := g.client.Login(context.Background(), &proto2.LoginRequest{
		Id:   loginID,
		User: *g.userID,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.NotFound {
				return nil, fmt.Errorf(FmtErrNotFound, err)
			}
		}
		return nil, fmt.Errorf(FmtErrInternalServer, err)
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
	_, err := g.client.DeleteLogin(context.Background(), &proto2.DeleteLoginRequest{
		Id:   loginID,
		User: *g.userID,
	})
	if err != nil {
		return fmt.Errorf(FmtErrInternalServer, err)
	}
	return nil
}

func (g *GRPCSender) AddText(text *models.Text) error {
	if g.userID == nil {
		return ErrAuthRequire
	}

	_, err := g.client.AddText(context.Background(), &proto2.AddTextRequest{
		Text: &proto2.Text{
			Id:       text.ID,
			Content:  text.Content,
			Metainfo: text.MetaInfo,
		},
		User: *g.userID,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.AlreadyExists {
				return fmt.Errorf(FmtErrAlreadyExists, err)
			}
		}
		return fmt.Errorf(FmtErrInternalServer, err)
	}
	return nil
}

func (g *GRPCSender) Text(textID string) (*models.Text, error) {
	if g.userID == nil {
		return nil, ErrAuthRequire
	}

	response, err := g.client.Text(context.Background(), &proto2.TextRequest{
		Id:   textID,
		User: *g.userID,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.NotFound {
				return nil, fmt.Errorf(FmtErrNotFound, err)
			}
		}
		return nil, fmt.Errorf(FmtErrInternalServer, err)
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
	_, err := g.client.DeleteText(context.Background(), &proto2.DeleteTextRequest{
		Id:   textID,
		User: *g.userID,
	})
	if err != nil {
		return fmt.Errorf(FmtErrInternalServer, err)
	}
	return nil
}

func (g *GRPCSender) AddBin(binary *models.Binary) error {
	if g.userID == nil {
		return ErrAuthRequire
	}

	_, err := g.client.AddBinary(context.Background(), &proto2.AddBinRequest{
		Binary: &proto2.Binary{
			Id:       binary.ID,
			Data:     binary.Data,
			Metainfo: binary.MetaInfo,
		},
		User: *g.userID,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.AlreadyExists {
				return fmt.Errorf(FmtErrAlreadyExists, err)
			}
		}
		return fmt.Errorf(FmtErrInternalServer, err)
	}
	return nil
}

func (g *GRPCSender) Bin(binID string) (*models.Binary, error) {
	if g.userID == nil {
		return nil, ErrAuthRequire
	}

	response, err := g.client.Binary(context.Background(), &proto2.BinRequest{
		Id:   binID,
		User: *g.userID,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.NotFound {
				return nil, fmt.Errorf(FmtErrNotFound, err)
			}
		}
		return nil, fmt.Errorf(FmtErrInternalServer, err)
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
	_, err := g.client.DeleteBinary(context.Background(), &proto2.DeleteBinRequest{
		Id:   binID,
		User: *g.userID,
	})
	if err != nil {
		return fmt.Errorf(FmtErrInternalServer, err)
	}
	return nil
}
