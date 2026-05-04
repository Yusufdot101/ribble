package grpc

import (
	"context"
	"errors"

	userpb "github.com/Yusufdot101/ripple-proto/golang/user/v4"
	"github.com/Yusufdot101/ripple/services/user/internal/application/core/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *Adapter) VerifyUsers(ctx context.Context, req *userpb.VerifyUsersRequest) (*userpb.VerifyUsersResponse, error) {
	allValid, err := a.asvc.VerifyUsers(ctx, req.Ids)
	if err != nil {
		var st error
		if errors.Is(err, domain.ErrInvalidID) {
			st = status.Error(codes.NotFound, err.Error())
		} else {
			st = status.Error(codes.Internal, err.Error())
		}
		return &userpb.VerifyUsersResponse{
			AllValid: false,
		}, st
	}

	return &userpb.VerifyUsersResponse{
		AllValid: allValid,
	}, nil
}

func (a *Adapter) GetUsers(ctx context.Context, req *userpb.GetUsersRequest) (*userpb.GetUsersResponse, error) {
	users, err := a.usvc.GetUsersByIDs(ctx, req.Ids)
	grpcUsers := []*userpb.User{}
	if err != nil {
		var st error
		if errors.Is(err, domain.ErrInvalidID) {
			st = status.Error(codes.NotFound, err.Error())
		} else {
			st = status.Error(codes.Internal, err.Error())
		}
		return &userpb.GetUsersResponse{
			Users: grpcUsers,
		}, st
	}

	for _, user := range users {
		grpcUsers = append(grpcUsers, &userpb.User{
			Name:     user.Name,
			Id:       uint32(user.ID),
			Sub:      user.Sub,
			Provider: user.Provider,
			Email:    user.Email,
		})
	}

	return &userpb.GetUsersResponse{
		Users: grpcUsers,
	}, nil
}

func (a *Adapter) SearchUsers(ctx context.Context, req *userpb.SearchUsersRequest) (*userpb.SearchUsersResponse, error) {
	users, err := a.usvc.SearchUsers(ctx, req.Query, req.Ids)
	grpcUsers := []*userpb.User{}
	if err != nil {
		return &userpb.SearchUsersResponse{Users: grpcUsers}, err
	}

	for _, user := range users {
		grpcUsers = append(grpcUsers, &userpb.User{
			Name:     user.Name,
			Email:    user.Email,
			Id:       uint32(user.ID),
			Sub:      user.Sub,
			Provider: user.Provider,
		})
	}

	return &userpb.SearchUsersResponse{Users: grpcUsers}, nil
}

func (a *Adapter) GetContacts(ctx context.Context, req *userpb.GetContactsRequest) (*userpb.GetContactsResponse, error) {
	users, err := a.usvc.GetContacts(ctx, req.Query, req.ExcludeIds, req.UserId)

	grpcUsers := []*userpb.User{}
	if err != nil {
		return &userpb.GetContactsResponse{Users: grpcUsers}, err
	}

	for _, user := range users {
		grpcUsers = append(grpcUsers, &userpb.User{
			Name:     user.Name,
			Id:       uint32(user.ID),
			Sub:      user.Sub,
			Provider: user.Provider,
			Email:    user.Email,
		})
	}
	return &userpb.GetContactsResponse{Users: grpcUsers}, nil
}
