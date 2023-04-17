package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/golang-jwt/jwt/v5"
)

const (
	tokenHeader = "Ark-Token"
	tokenIssuer = "ark-client"
)

var errNoToken = errors.New("no token provided")

type Interceptor struct {
	key []byte
}

func NewInterceptor(signingKey string) *Interceptor {
	return &Interceptor{key: []byte(signingKey)}
}

func (i *Interceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		if req.Spec().IsClient {
			claims := &jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
				Issuer:    tokenIssuer,
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			signedToken, err := token.SignedString(i.key)
			if err != nil {
				fmt.Println("error signing token", err)
				return nil, err
			}

			fmt.Println("unary token", signedToken)
		} else if req.Header().Get(tokenHeader) == "" {
			tokenString := req.Header().Get(tokenHeader)
			if tokenString == "" {
				return nil, connect.NewError(connect.CodeUnauthenticated, errNoToken)
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("unexpected signing method: %v", token.Header["alg"]))
				}

				expiry, err := token.Claims.GetExpirationTime()
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, errors.New("error getting expiration time claim"))
				}

				if expiry.Before(time.Now()) {
					return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("token expired"))
				}

				issuer, err := token.Claims.GetIssuer()
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, errors.New("error getting issuer claim"))
				}

				if issuer != tokenIssuer {
					return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid issuer"))
				}

				return i.key, nil
			})
			if err != nil {
				return nil, err
			}

			if !token.Valid {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid token"))
			}
		}
		return next(ctx, req)
	}
}

func (i *Interceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return func(ctx context.Context, spec connect.Spec) connect.StreamingClientConn {
		conn := next(ctx, spec)
		claims := &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer:    tokenIssuer,
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := token.SignedString(i.key)
		if err != nil {
			fmt.Println("error signing token", err)
			return conn
		}

		fmt.Println("token", signedToken)

		conn.RequestHeader().Set(tokenHeader, signedToken)

		return conn
	}
}

func (i *Interceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		tokenString := conn.RequestHeader().Get(tokenHeader)
		if tokenString == "" {
			return connect.NewError(connect.CodeUnauthenticated, errNoToken)
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("unexpected signing method: %v", token.Header["alg"]))
			}

			expiry, err := token.Claims.GetExpirationTime()
			if err != nil {
				return nil, connect.NewError(connect.CodeInternal, errors.New("error getting expiration time claim"))
			}

			if expiry.Before(time.Now()) {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("token expired"))
			}

			issuer, err := token.Claims.GetIssuer()
			if err != nil {
				return nil, connect.NewError(connect.CodeInternal, errors.New("error getting issuer claim"))
			}

			if issuer != tokenIssuer {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid issuer"))
			}

			return i.key, nil
		})
		if err != nil {
			return err
		}

		if !token.Valid {
			return connect.NewError(connect.CodeUnauthenticated, errors.New("invalid token"))
		}

		return next(ctx, conn)
	}
}