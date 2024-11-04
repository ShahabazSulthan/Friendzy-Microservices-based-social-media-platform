package handler_post

import (
	"context"
	"fmt"
	"time"

	responsemodel_post "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/post_relation_service/model/responsemodel"
	"github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/post_relation_service/pb"
	"github.com/gofiber/fiber/v2"
)

type RelationHandler struct {
	Client pb.PostNrelServiceClient
}

func NewRelationHandler(client *pb.PostNrelServiceClient) *RelationHandler {
	return &RelationHandler{Client: *client}
}

func (svc *RelationHandler) Follow(ctx *fiber.Ctx) error {
	userid := ctx.Locals("userId")
	userId := fmt.Sprint(userid)

	userBId := ctx.Params("followingid")

	if userBId == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "failed request(possible-reason:no input)",
				Error:      "no userBId (\":followingid\") param found in request.",
			})
	}

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := svc.Client.Follow(context, &pb.RequestFollowUnFollow{
		UserId:  userId,
		UserBId: userBId,
	})

	if err != nil {
		fmt.Println("----------postNrel service down--------")

		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "failed to follow userB",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "failed to follow userB",
				Data:       resp,
				Error:      resp.ErrorMessage,
			})
	}

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_post.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "followed succesfully",
			Data:       resp,
			Error:      nil,
		})

}

func (svc *RelationHandler) UnFollow(ctx *fiber.Ctx) error {

	userid := ctx.Locals("userId")
	userId := fmt.Sprint(userid)

	userBId := ctx.Params("unfollowingid")

	if userBId == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "failed request(possible-reason:no input)",
				Error:      "no userBId (\":unfollowingid\") param found in request.",
			})
	}

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := svc.Client.UnFollow(context, &pb.RequestFollowUnFollow{
		UserId:  userId,
		UserBId: userBId,
	})

	if err != nil {
		fmt.Println("----------postNrel service down--------")

		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "failed to unfollow userB",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "failed to unfollow userB",
				Data:       resp,
				Error:      resp.ErrorMessage,
			})
	}

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_post.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "unfollowed succesfully",
			Data:       resp,
			Error:      nil,
		})
}

