package interceptor

import (
	"context"
	"encoding/hex"
	"strings"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"gophkeeper/internal/client/crypto"
)

type AuthClient struct {
	rsa *crypto.UserSession
}

func NewAuthClient(rsa *crypto.UserSession) *AuthClient {
	return &AuthClient{rsa: rsa}
}

func (client *AuthClient) Unary() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if strings.Contains(method, "NewSessionID") {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		userSign, err := client.rsa.UserSign()
		if err != nil {
			log.Error().Err(err).Msg("RegisterUser EncryptOAEP signing error")
			return err
		}
		ctx = metadata.AppendToOutgoingContext(ctx, "userSession", client.rsa.GetSessionID())
		ctx = metadata.AppendToOutgoingContext(ctx, "userSign", hex.EncodeToString(userSign))
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
