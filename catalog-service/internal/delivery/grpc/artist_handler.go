package v1

import (
	"context"

	"github.com/Lama189/soundwave-platform/catalog-service/internal/service/artist"
	catalogpb "github.com/Lama189/soundwave-platform/gen/go/catalog/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ArtistGRPCServer struct {
	catalogpb.UnimplementedArtistServiceServer
	service ArtistService
}

func NewArtistGRPCServer(service ArtistService) *ArtistGRPCServer {
	return &ArtistGRPCServer{
		service: service,
	}
}

func (s *ArtistGRPCServer) CreateArtist(ctx context.Context, req *catalogpb.CreateArtistRequest) (*catalogpb.CreateArtistResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id format")
	}

	dto := artist.CreateArtistDTO{
		Name: req.GetName(),
		Bio:  req.GetBio(),
	}

	a, err := s.service.Create(ctx, userID, dto)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &catalogpb.CreateArtistResponse{
		Artist: toProtoArtist(a),
	}, nil
}

func (s *ArtistGRPCServer) GetArtist(ctx context.Context, req *catalogpb.GetArtistRequest) (*catalogpb.GetArtistResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid artist id format")
	}

	a, err := s.service.GetByID(ctx, id)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &catalogpb.GetArtistResponse{
		Artist: toProtoArtist(a),
	}, nil
}

func (s *ArtistGRPCServer) GetArtistByUserId(ctx context.Context, req *catalogpb.GetArtistByUserIdRequest) (*catalogpb.GetArtistByUserIdResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id format")
	}

	a, err := s.service.GetByUserID(ctx, userID)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &catalogpb.GetArtistByUserIdResponse{
		Artist: toProtoArtist(a),
	}, nil
}

func (s *ArtistGRPCServer) UpdateArtist(ctx context.Context, req *catalogpb.UpdateArtistRequest) (*catalogpb.UpdateArtistResponse, error) {
	artistID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid artist id format")
	}

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id format")
	}

	dto := artist.UpdateArtistDTO{
		Name: req.GetName(),
		Bio:  req.GetBio(),
	}

	a, err := s.service.Update(ctx, userID, artistID, dto)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &catalogpb.UpdateArtistResponse{
		Artist: toProtoArtist(a),
	}, nil
}

func (s *ArtistGRPCServer) DeleteArtist(ctx context.Context, req *catalogpb.DeleteArtistRequest) (*catalogpb.DeleteArtistResponse, error) {
	artistID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid artist id format")
	}

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id format")
	}

	if err := s.service.Delete(ctx, userID, artistID); err != nil {
		return nil, mapDomainError(err)
	}

	return &catalogpb.DeleteArtistResponse{}, nil
}

func toProtoArtist(a *artist.ArtistOutputDTO) *catalogpb.Artist {
	if a == nil {
		return nil
	}
	return &catalogpb.Artist{
		Id:        a.ID.String(),
		UserId:    a.UserID.String(),
		Name:      a.Name,
		Bio:       a.Bio,
		CreatedAt: timestamppb.New(a.CreatedAt),
		UpdatedAt: timestamppb.New(a.UpdatedAt),
	}
}
