// Package api. Реализует UnimplementedGophKeeperServerServer, те все серверное взаимодействие
package api

import (
	"context"
	"crypto/sha1"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/ncyellow/GophKeeper/internal/models"
	proto2 "github.com/ncyellow/GophKeeper/internal/proto"
	"github.com/ncyellow/GophKeeper/internal/server/config"
	"github.com/ncyellow/GophKeeper/internal/server/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCServer структура для реализации grpc сервера. Для https сервера я делал полное тестирование.
// Потому тут я повторяться не буду и второй раз тестировать тоже самое не буду
type GRPCServer struct {
	proto2.UnimplementedGophKeeperServerServer
	conf *config.Config
	repo storage.Storage
}

// NewServer конструктор
func NewServer(repo storage.Storage, conf *config.Config) *GRPCServer {
	return &GRPCServer{
		repo: repo,
		conf: conf,
	}
}

// Register регистрация пользователя
func (s *GRPCServer) Register(ctx context.Context, req *proto2.RegisterRequest) (*proto2.RegisterResponse, error) {
	login := req.GetLogin()
	originalPassword := req.GetPassword()

	pwd := sha1.New()
	pwd.Write([]byte(originalPassword))
	hashPwd := fmt.Sprintf("%x", pwd.Sum(nil))

	user := models.User{
		Login:    login,
		Password: hashPwd,
	}

	// Выполняем попытку регистрации
	_, err := s.repo.Register(ctx, user)
	if err != nil {
		return nil, status.Error(codes.AlreadyExists, "")
	}

	repoUser, err := s.repo.User(ctx, user.Login, user.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "")
	}

	return &proto2.RegisterResponse{
		User: repoUser.UserID,
	}, nil
}

// SignIn аутентификация
func (s *GRPCServer) SignIn(ctx context.Context, req *proto2.RegisterRequest) (*proto2.RegisterResponse, error) {
	login := req.GetLogin()
	originalPassword := req.GetPassword()

	pwd := sha1.New()
	pwd.Write([]byte(originalPassword))
	hashPwd := fmt.Sprintf("%x", pwd.Sum(nil))

	user := models.User{
		Login:    login,
		Password: hashPwd,
	}

	repoUser, err := s.repo.User(ctx, user.Login, user.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "")
	}

	return &proto2.RegisterResponse{
		User: repoUser.UserID,
	}, nil
}

// AddCard зарегистрировать новую карту
func (s *GRPCServer) AddCard(ctx context.Context, req *proto2.AddCardRequest) (*proto2.AddCardResponse, error) {
	var response proto2.AddCardResponse
	card := *req.GetCard()
	userID := req.GetUser()
	err := s.repo.AddCard(ctx, userID, models.Card{
		ID:       card.GetId(),
		FIO:      card.GetFio(),
		Number:   card.GetFio(),
		Date:     card.GetDate(),
		CVV:      card.GetCvv(),
		MetaInfo: card.GetMetainfo(),
	})
	if err != nil {
		return nil, status.Error(codes.AlreadyExists, "")
	}
	return &response, nil
}

// AddLogin зарегистрировать новый логин
func (s *GRPCServer) AddLogin(ctx context.Context, req *proto2.AddLoginRequest) (*proto2.AddLoginResponse, error) {
	var response proto2.AddLoginResponse
	login := *req.GetLogin()
	userID := req.GetUser()
	err := s.repo.AddLogin(ctx, userID, models.Login{
		ID:       login.GetId(),
		Login:    login.GetLogin(),
		Password: login.GetPassword(),
		MetaInfo: login.GetMetainfo(),
	})
	if err != nil {
		return nil, status.Error(codes.AlreadyExists, "")
	}
	return &response, nil
}

// AddText зарегистрировать новый текст по пользователю и идентификатору
func (s *GRPCServer) AddText(ctx context.Context, req *proto2.AddTextRequest) (*proto2.AddTextResponse, error) {
	var response proto2.AddTextResponse
	text := *req.GetText()
	userID := req.GetUser()
	err := s.repo.AddText(ctx, userID, models.Text{
		ID:       text.GetId(),
		Content:  text.GetContent(),
		MetaInfo: text.GetMetainfo(),
	})
	if err != nil {
		return nil, status.Error(codes.AlreadyExists, "")
	}
	return &response, nil
}

// AddBinary зарегистрировать бинарные данные по пользователю и идентификатору
func (s *GRPCServer) AddBinary(ctx context.Context, req *proto2.AddBinRequest) (*proto2.AddBinResponse, error) {
	var response proto2.AddBinResponse
	text := *req.GetBinary()
	userID := req.GetUser()
	err := s.repo.AddBinary(ctx, userID, models.Binary{
		ID:       text.GetId(),
		Data:     text.GetData(),
		MetaInfo: text.GetMetainfo(),
	})
	if err != nil {
		return nil, status.Error(codes.AlreadyExists, "")
	}
	return &response, nil
}

// Card вернуть данные конкретной карты
func (s *GRPCServer) Card(ctx context.Context, req *proto2.CardRequest) (*proto2.CardResponse, error) {
	cardID := req.GetId()
	userID := req.GetUser()
	card, err := s.repo.Card(ctx, userID, cardID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, status.Error(codes.NotFound, "")
		}
		return nil, status.Error(codes.Internal, "")

	}
	return &proto2.CardResponse{
		Card: &proto2.Card{
			Id:       card.ID,
			Fio:      card.FIO,
			Number:   card.Number,
			Date:     card.Date,
			Cvv:      card.CVV,
			Metainfo: card.MetaInfo,
		},
	}, nil
}

// Login вернуть данные по конкретному логину
func (s *GRPCServer) Login(ctx context.Context, req *proto2.LoginRequest) (*proto2.LoginResponse, error) {
	loginID := req.GetId()
	userID := req.GetUser()
	login, err := s.repo.Login(ctx, userID, loginID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, status.Error(codes.NotFound, "")
		}
		return nil, status.Error(codes.Internal, "")

	}
	return &proto2.LoginResponse{
		Login: &proto2.Login{
			Id:       login.ID,
			Login:    login.Login,
			Password: login.Password,
			Metainfo: login.MetaInfo,
		},
	}, nil
}

// Text вернуть данные по конкретному тексту
func (s *GRPCServer) Text(ctx context.Context, req *proto2.TextRequest) (*proto2.TextResponse, error) {
	textID := req.GetId()
	userID := req.GetUser()
	text, err := s.repo.Text(ctx, userID, textID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, status.Error(codes.NotFound, "")
		}
		return nil, status.Error(codes.Internal, "")

	}
	return &proto2.TextResponse{
		Text: &proto2.Text{
			Id:       text.ID,
			Content:  text.Content,
			Metainfo: text.MetaInfo,
		},
	}, nil
}

// Binary вернуть данные конкретных бинарных данных
func (s *GRPCServer) Binary(ctx context.Context, req *proto2.BinRequest) (*proto2.BinResponse, error) {
	binID := req.GetId()
	userID := req.GetUser()
	bin, err := s.repo.Binary(ctx, userID, binID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, status.Error(codes.NotFound, "")
		}
		return nil, status.Error(codes.Internal, "")

	}
	return &proto2.BinResponse{
		Binary: &proto2.Binary{
			Id:       bin.ID,
			Data:     bin.Data,
			Metainfo: bin.MetaInfo,
		},
	}, nil
}

// DeleteCard удалить карту по пользователю и идентификатору
func (s *GRPCServer) DeleteCard(ctx context.Context, req *proto2.DeleteCardRequest) (*proto2.DeleteCardResponse, error) {
	var response proto2.DeleteCardResponse
	dataID := req.GetId()
	userID := req.GetUser()
	err := s.repo.DeleteCard(ctx, userID, dataID)
	if err != nil {
		return nil, status.Error(codes.Internal, "")
	}
	return &response, nil
}

// DeleteLogin удалить логин по пользователю и идентификатору
func (s *GRPCServer) DeleteLogin(ctx context.Context, req *proto2.DeleteLoginRequest) (*proto2.DeleteLoginResponse, error) {
	var response proto2.DeleteLoginResponse
	dataID := req.GetId()
	userID := req.GetUser()
	err := s.repo.DeleteLogin(ctx, userID, dataID)
	if err != nil {
		return nil, status.Error(codes.Internal, "")
	}
	return &response, nil
}

// DeleteText текст по пользователю и идентификатору
func (s *GRPCServer) DeleteText(ctx context.Context, req *proto2.DeleteTextRequest) (*proto2.DeleteTextResponse, error) {
	var response proto2.DeleteTextResponse
	dataID := req.GetId()
	userID := req.GetUser()
	err := s.repo.DeleteText(ctx, userID, dataID)
	if err != nil {
		return nil, status.Error(codes.Internal, "")
	}
	return &response, nil
}

// DeleteBinary удалить бинарные данные по пользователю и идентификатору
func (s *GRPCServer) DeleteBinary(ctx context.Context, req *proto2.DeleteBinRequest) (*proto2.DeleteBinResponse, error) {
	var response proto2.DeleteBinResponse
	dataID := req.GetId()
	userID := req.GetUser()
	err := s.repo.DeleteBinary(ctx, userID, dataID)
	if err != nil {
		return nil, status.Error(codes.Internal, "")
	}
	return &response, nil
}
