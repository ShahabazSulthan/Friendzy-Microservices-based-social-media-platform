package handler_post

import (
	"context"
	"fmt"
	"log"
	"time"

	byteconverter "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Utils/byteConverter"
	mediafileformatchecker "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Utils/mediaFileFormatChecker"
	requestmodel_post "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/post_relation_service/model/requestmodel"
	responsemodel_post "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/post_relation_service/model/responsemodel"
	"github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/post_relation_service/pb"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type PostHandler struct {
	Client pb.PostNrelServiceClient
}

func NewPostHandler(client *pb.PostNrelServiceClient) *PostHandler {
	return &PostHandler{Client: *client}
}

func (svc *PostHandler) AddNewPost(ctx *fiber.Ctx) error {

	var postData requestmodel_post.AddPostData
	var respPostData responsemodel_post.AddPostResp

	userId := ctx.Locals("userId")
	postData.UserId = fmt.Sprint(userId)

	if err := ctx.BodyParser(&postData); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "can't add post(possible-reason:no json input)",
				Error:      err.Error(),
			})
	}

	//fiber's ctx.BodyParser can't parse files(*multipart.FileHeader),
	//so we have to manually access the Multipart form and read the files form it.
	form, err := ctx.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["media"]
	postData.Media = append(postData.Media, files...)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(postData)
	if err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			for _, e := range ve {
				switch e.Field() {
				case "Caption":
					respPostData.Caption = "should contain less than 60 letters"
				case "UserId":
					respPostData.UserId = "No userId got"
				case "Media":
					respPostData.Media = "you can't add a post without a image/video"
				}
			}
		}
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "can't add post",
				Data:       respPostData,
				Error:      err.Error(),
			})
	}

	numFiles := len(postData.Media)
	if numFiles < 1 || numFiles > 5 {
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "can't add post",
				Data:       nil,
				Error:      "you can only add 5 image/video in a post",
			})
	}

	for _, media := range postData.Media {
		if media.Size > 5*1024*1024 { // 5 MB limit
			return ctx.Status(fiber.ErrBadRequest.Code).
				JSON(responsemodel_post.CommonResponse{
					StatusCode: fiber.ErrBadRequest.Code,
					Message:    "can't add post",
					Data:       nil,
					Error:      "yfile size exceeds the limit (5MB)",
				})
		}
	}

	var mediaData []*pb.SingleMedia
	for _, fileHeader := range postData.Media {
		file, err := fileHeader.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		contentType, err := mediafileformatchecker.Mediafileformatchecker(file)
		if err != nil {
			return ctx.Status(fiber.ErrBadRequest.Code).
				JSON(responsemodel_post.CommonResponse{
					StatusCode: fiber.ErrBadRequest.Code,
					Message:    "can't add post",
					Data:       nil,
					Error:      err.Error(),
				})
		}

		content, err := byteconverter.MultipartFileheaderToBytes(&file)
		if err != nil {
			fmt.Println("-------------byteconverter-down---------")
		}

		mediaData = append(mediaData, &pb.SingleMedia{Media: content, ContentType: *contentType})

	}

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := svc.Client.AddNewPost(context, &pb.RequestAddPost{
		UserId:  postData.UserId,
		Caption: postData.Caption,
		Media:   mediaData,
	})

	if err != nil {
		fmt.Println("----------postNrel service down--------")

		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "can't add post",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "can't add post",
				Data:       resp,
				Error:      resp.ErrorMessage,
			})
	}

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_post.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "Post added succesfully",
			Data:       resp,
			Error:      nil,
		})
}

func (svc *PostHandler) GetAllPostByUser(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId")
	sendingId := ctx.Query("userbid", fmt.Sprint(userId))

	limit, offset := ctx.Query("limit", "12"), ctx.Query("offset", "0")

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := svc.Client.GetAllPostByUser(context, &pb.RequestGetAllPosts{
		UserId: sendingId,
		Limit:  limit,
		OffSet: offset,
	})

	if err != nil {
		fmt.Println("----------postNrel service down--------")

		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "can't fetch Posts",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "can't fetch Posts",
				Data:       resp,
				Error:      resp.ErrorMessage,
			})
	}

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_post.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "Posts fetched succesfully",
			Data:       resp,
			Error:      nil,
		})
}

func (svc *PostHandler) DeletePost(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId")
	postId := ctx.Params("postid")

	if fmt.Sprint(userId) == "" || postId == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "can't delete post",
				Data:       nil,
				Error:      "no postid found in request",
			})
	}

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := svc.Client.DeletePost(context, &pb.RequestDeletePost{
		UserId: fmt.Sprint(userId),
		PostId: postId,
	})

	if err != nil {
		fmt.Println("----------postNrel service down--------")

		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "can't delete Posts",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "can't delete Posts",
				Data:       resp,
				Error:      resp.ErrorMessage,
			})
	}

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_post.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "Post deleted succesfully",
		})
}

func (svc *PostHandler) EditPost(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId")

	var editInput requestmodel_post.EditPost
	var respPostEdit responsemodel_post.EditPostResp

	editInput.UserId = fmt.Sprint(userId)

	if err := ctx.BodyParser(&editInput); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "can't add post(possible-reason:no json input)",
				Error:      err.Error(),
			})
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(editInput)
	if err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			for _, e := range ve {
				switch e.Field() {
				case "Caption":
					respPostEdit.Caption = "should contain less than 60 letters"
				case "UserId":
					respPostEdit.UserId = "No userId got from header"
				case "PostId":
					respPostEdit.PostId = "no postid found in request"

				}
			}
			return ctx.Status(fiber.ErrBadRequest.Code).
				JSON(responsemodel_post.CommonResponse{
					StatusCode: fiber.ErrBadRequest.Code,
					Message:    "can't edit post",
					Data:       respPostEdit,
					Error:      err.Error(),
				})
		}
	}

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := svc.Client.EditPost(context, &pb.RequestEditPost{
		UserId:  editInput.UserId,
		PostId:  editInput.PostId,
		Caption: editInput.Caption,
	})

	if err != nil {
		fmt.Println("----------postNrel service down--------")

		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "can't edit Posts",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "can't edit Posts",
				Data:       resp,
				Error:      resp.ErrorMessage,
			})
	}

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_post.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "Posts edited succesfully",
		})

}

func (svc *PostHandler) LikePost(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId")
	postId := ctx.Params("postid")

	if fmt.Sprint(userId) == "" || postId == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "can't like post",
				Data:       nil,
				Error:      "no postid found in request",
			})
	}

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := svc.Client.LikePost(context, &pb.RequestLikeUnlikePost{
		UserId: fmt.Sprint(userId),
		PostId: postId,
	})

	if err != nil {
		fmt.Println("----------postNrel service down--------")

		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "can't like Post",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "can't like Post",
				Data:       resp,
				Error:      resp.ErrorMessage,
			})
	}

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_post.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "Post liked succesfully",
		})
}

func (svc *PostHandler) UnLikePost(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId")
	postId := ctx.Params("postid")

	if fmt.Sprint(userId) == "" || postId == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "can't like post",
				Data:       nil,
				Error:      "no postid found in request",
			})
	}

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := svc.Client.UnLikePost(context, &pb.RequestLikeUnlikePost{
		UserId: fmt.Sprint(userId),
		PostId: postId,
	})

	if err != nil {
		fmt.Println("----------postNrel service down--------")

		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "can't unlike Post",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "can't unlike Post",
				Data:       resp,
				Error:      resp.ErrorMessage,
			})
	}

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_post.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "Post unliked succesfully",
		})

}
func (svc *PostHandler) GetMostLovedPostsFromGlobalUser(ctx *fiber.Ctx) error {

	userId := ctx.Locals("userId")
	limit, offset := ctx.Query("limit", "22"), ctx.Query("offset", "0")

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := svc.Client.GetMostLovedPostsFromGlobalUser(context, &pb.RequestGetAllPosts{
		UserId: fmt.Sprint(userId),
		Limit:  limit,
		OffSet: offset,
	})

	if err != nil {
		fmt.Println("----------postNrel service down--------")

		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "can't fetch Posts",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "can't fetch Posts",
				Data:       resp,
				Error:      resp.ErrorMessage,
			})
	}

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_post.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "Posts fetched succesfully",
			Data:       resp,
			Error:      nil,
		})
}

func (svc *PostHandler) GetAllRelatedPostsForHomeScreen(ctx *fiber.Ctx) error {

	userId := ctx.Locals("userId")
	limit, offset := ctx.Query("limit", "22"), ctx.Query("offset", "0")

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := svc.Client.GetAllRelatedPostsForHomeScreen(context, &pb.RequestGetAllPosts{
		UserId: fmt.Sprint(userId),
		Limit:  limit,
		OffSet: offset,
	})

	if err != nil {
		fmt.Println("----------postNrel service down--------")

		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "can't fetch Posts",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "can't fetch Posts",
				Data:       resp,
				Error:      resp.ErrorMessage,
			})
	}

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_post.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "Posts fetched succesfully",
			Data:       resp,
			Error:      nil,
		})
}

func (svc *PostHandler) GetAllRandomPosts(ctx *fiber.Ctx) error {

	limit, offset := ctx.Query("limit", "22"), ctx.Query("offset", "0")

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := svc.Client.GetRandomPosts(context, &pb.RequestGetRandomPosts{
		Limit:  limit,
		OffSet: offset,
	})

	if err != nil {
		fmt.Println("----------postNrel service down--------")

		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "can't fetch Posts",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "can't fetch Posts",
				Data:       resp,
				Error:      resp.ErrorMessage,
			})
	}

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_post.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "Posts fetched succesfully",
			Data:       resp,
			Error:      nil,
		})
}
