package interceptor

import (
	"context"
	"encoding/hex"
	"strings"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"gophkeeper/internal/server/crypto"
)

type AuthInterceptor struct {
	rsa *crypto.Sessions
}

func NewAuthInterceptor(rsa *crypto.Sessions) *AuthInterceptor {
	return &AuthInterceptor{rsa: rsa}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Error().Msgf("interceptor receieved method = %s", info.FullMethod)
		if strings.Contains(info.FullMethod, "NewSessionID") {
			return handler(ctx, req)
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
		}
		session := md.Get("userSession")
		sign := md.Get("userSign")
		userSign, _ := hex.DecodeString(sign[0])
		userID, err := interceptor.rsa.CheckSign(session[0], userSign)
		if err != nil {
			log.Error().Err(err).Msg("UserData CheckSign error")
			return nil, status.Error(codes.Unauthenticated, "incorrect sign encryption")
		}
		if strings.Contains(info.FullMethod, "NewUser") || strings.Contains(info.FullMethod, "LoginUser") {
			return handler(ctx, req)
		}
		if userID == "" {
			log.Error().Err(err).Msg("UserData userID empty")
			return nil, status.Error(codes.Unauthenticated, "incorrect sign encryption")
		}

		return handler(ctx, req)
	}
}
