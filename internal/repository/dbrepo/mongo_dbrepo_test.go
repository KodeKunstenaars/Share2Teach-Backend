package dbrepo

import (
	"backend/internal/models"
	"backend/pkg/db"
	"context"
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	testUserJoe = models.User{
		ID:            primitive.NewObjectID(),
		FirstName:     "John",
		LastName:      "Doe",
		Email:         "john.doe@mail.com",
		Password:      "123",
		Role:          "admin",
		Qualification: "M.Sc.",
	}
	testDocument = models.Document{
		ID:        primitive.NewObjectID(),
		Title:     "Chapter 1 - Introduction to Go",
		Subject:   "Computer Science",
		Grade:     "12",
		CreatedAt: time.Now(),
		UserID:    testUserJoe.ID,
		Moderated: false,
		Reported:  false,
		RatingID:  primitive.NewObjectID(),
	}
	testRating = models.Rating{
		ID:            primitive.NewObjectID(),
		DocID:         testDocument.ID,
		TimesRated:    2,
		TotalRating:   10,
		AverageRating: 5,
	}
)

func TestMongoDBRepo_ChangeUserPassword(t *testing.T) {
	type fields struct {
		userInfoCollection      db.Collection
		passwordResetCollection db.Collection
		metadataCollection      db.Collection
		ratingsCollection       db.Collection
		faqsCollection          db.Collection
		moderateCollection      db.Collection
		reportsCollection       db.Collection
	}
	type args struct {
		userID      primitive.ObjectID
		newPassword string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MongoDBRepo{
				userInfoCollection:      tt.fields.userInfoCollection,
				passwordResetCollection: tt.fields.passwordResetCollection,
				metadataCollection:      tt.fields.metadataCollection,
				ratingsCollection:       tt.fields.ratingsCollection,
				faqsCollection:          tt.fields.faqsCollection,
				moderateCollection:      tt.fields.moderateCollection,
				reportsCollection:       tt.fields.reportsCollection,
			}
			if err := m.ChangeUserPassword(tt.args.userID, tt.args.newPassword); (err != nil) != tt.wantErr {
				t.Errorf("ChangeUserPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMongoDBRepo_CreateDocumentRating(t *testing.T) {
	type fields struct {
		userInfoCollection      db.Collection
		passwordResetCollection db.Collection
		metadataCollection      db.Collection
		ratingsCollection       db.Collection
		faqsCollection          db.Collection
		moderateCollection      db.Collection
		reportsCollection       db.Collection
	}
	type args struct {
		initialRating *models.Rating
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MongoDBRepo{
				userInfoCollection:      tt.fields.userInfoCollection,
				passwordResetCollection: tt.fields.passwordResetCollection,
				metadataCollection:      tt.fields.metadataCollection,
				ratingsCollection:       tt.fields.ratingsCollection,
				faqsCollection:          tt.fields.faqsCollection,
				moderateCollection:      tt.fields.moderateCollection,
				reportsCollection:       tt.fields.reportsCollection,
			}
			if err := m.CreateDocumentRating(tt.args.initialRating); (err != nil) != tt.wantErr {
				t.Errorf("CreateDocumentRating() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMongoDBRepo_FindDocuments(t *testing.T) {
	type fields struct {
		userInfoCollection      db.Collection
		passwordResetCollection db.Collection
		metadataCollection      db.Collection
		ratingsCollection       db.Collection
		faqsCollection          db.Collection
		moderateCollection      db.Collection
		reportsCollection       db.Collection
	}
	type args struct {
		title       string
		subject     string
		grade       string
		correctRole bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Document
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MongoDBRepo{
				userInfoCollection:      tt.fields.userInfoCollection,
				passwordResetCollection: tt.fields.passwordResetCollection,
				metadataCollection:      tt.fields.metadataCollection,
				ratingsCollection:       tt.fields.ratingsCollection,
				faqsCollection:          tt.fields.faqsCollection,
				moderateCollection:      tt.fields.moderateCollection,
				reportsCollection:       tt.fields.reportsCollection,
			}
			got, err := m.FindDocuments(tt.args.title, tt.args.subject, tt.args.grade, tt.args.correctRole)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindDocuments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindDocuments() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMongoDBRepo_GetDocumentByID(t *testing.T) {
	type fields struct {
		userInfoCollection      db.Collection
		passwordResetCollection db.Collection
		metadataCollection      db.Collection
		ratingsCollection       db.Collection
		faqsCollection          db.Collection
		moderateCollection      db.Collection
		reportsCollection       db.Collection
	}
	type args struct {
		id primitive.ObjectID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Document
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MongoDBRepo{
				userInfoCollection:      tt.fields.userInfoCollection,
				passwordResetCollection: tt.fields.passwordResetCollection,
				metadataCollection:      tt.fields.metadataCollection,
				ratingsCollection:       tt.fields.ratingsCollection,
				faqsCollection:          tt.fields.faqsCollection,
				moderateCollection:      tt.fields.moderateCollection,
				reportsCollection:       tt.fields.reportsCollection,
			}
			got, err := m.GetDocumentByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDocumentByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDocumentByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMongoDBRepo_GetDocumentRating(t *testing.T) {
	type fields struct {
		userInfoCollection      db.Collection
		passwordResetCollection db.Collection
		metadataCollection      db.Collection
		ratingsCollection       db.Collection
		faqsCollection          db.Collection
		moderateCollection      db.Collection
		reportsCollection       db.Collection
	}
	type args struct {
		docID primitive.ObjectID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Rating
		wantErr bool
	}{
		{
			name: "get valid rating",
			args: args{docID: testDocument.ID},
			fields: fields{
				ratingsCollection: &db.MongoCollectionMock{
					FindOneFunc: func(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
						return mongo.NewSingleResultFromDocument(testRating, nil, nil)
					},
				},
			},
			want:    &testRating,
			wantErr: false,
		},
		{
			name: "rating not found",
			args: args{docID: testDocument.ID},
			fields: fields{
				ratingsCollection: &db.MongoCollectionMock{
					FindOneFunc: func(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
						return mongo.NewSingleResultFromDocument(nil, nil, nil)
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MongoDBRepo{
				userInfoCollection:      tt.fields.userInfoCollection,
				passwordResetCollection: tt.fields.passwordResetCollection,
				metadataCollection:      tt.fields.metadataCollection,
				ratingsCollection:       tt.fields.ratingsCollection,
				faqsCollection:          tt.fields.faqsCollection,
				moderateCollection:      tt.fields.moderateCollection,
				reportsCollection:       tt.fields.reportsCollection,
			}
			got, err := m.GetDocumentRating(tt.args.docID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDocumentRating() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDocumentRating() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMongoDBRepo_GetFAQs(t *testing.T) {
	type fields struct {
		userInfoCollection      db.Collection
		passwordResetCollection db.Collection
		metadataCollection      db.Collection
		ratingsCollection       db.Collection
		faqsCollection          db.Collection
		moderateCollection      db.Collection
		reportsCollection       db.Collection
	}
	tests := []struct {
		name    string
		fields  fields
		want    []models.FAQs
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MongoDBRepo{
				userInfoCollection:      tt.fields.userInfoCollection,
				passwordResetCollection: tt.fields.passwordResetCollection,
				metadataCollection:      tt.fields.metadataCollection,
				ratingsCollection:       tt.fields.ratingsCollection,
				faqsCollection:          tt.fields.faqsCollection,
				moderateCollection:      tt.fields.moderateCollection,
				reportsCollection:       tt.fields.reportsCollection,
			}
			got, err := m.GetFAQs()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFAQs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFAQs() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMongoDBRepo_GetUserByEmail(t *testing.T) {
	type fields struct {
		userInfoCollection      db.Collection
		passwordResetCollection db.Collection
		metadataCollection      db.Collection
		ratingsCollection       db.Collection
		faqsCollection          db.Collection
		moderateCollection      db.Collection
		reportsCollection       db.Collection
	}
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.User
		wantErr bool
	}{
		{
			name: "valid email",
			args: args{
				email: "john.doe@mail.com",
			},
			fields: fields{
				userInfoCollection: &db.MongoCollectionMock{
					FindOneFunc: func(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {

						// make sure the filter contains the email
						if filter == nil {
							return mongo.NewSingleResultFromDocument(nil, nil, nil)
						}

						// make sure the filter is a bson.M
						filterMap, ok := filter.(bson.M)
						if !ok {
							return mongo.NewSingleResultFromDocument(nil, nil, nil)
						}

						// make sure the email is in the filter
						email, ok := filterMap["email"]
						if !ok {
							return mongo.NewSingleResultFromDocument(nil, nil, nil)
						}

						// make sure the email is a string
						_, ok = email.(string)
						if !ok {
							return mongo.NewSingleResultFromDocument(nil, nil, nil)
						}

						return mongo.NewSingleResultFromDocument(models.User{
							ID:            primitive.ObjectID{},
							FirstName:     "John",
							LastName:      "Doe",
							Email:         "john.doe@mail.com",
							Password:      "123",
							Role:          "admin",
							Qualification: "M.Sc.",
						}, nil, nil)
					},
				},
			},
			wantErr: false,
			want: models.User{
				ID:            primitive.ObjectID{},
				FirstName:     "John",
				LastName:      "Doe",
				Email:         "john.doe@mail.com",
				Password:      "123",
				Role:          "admin",
				Qualification: "M.Sc.",
			},
		},
		{
			name: "error finding one",
			fields: fields{
				userInfoCollection: &db.MongoCollectionMock{
					FindOneFunc: func(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
						return mongo.NewSingleResultFromDocument(nil, nil, nil)
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MongoDBRepo{
				userInfoCollection:      tt.fields.userInfoCollection,
				passwordResetCollection: tt.fields.passwordResetCollection,
				metadataCollection:      tt.fields.metadataCollection,
				ratingsCollection:       tt.fields.ratingsCollection,
				faqsCollection:          tt.fields.faqsCollection,
				moderateCollection:      tt.fields.moderateCollection,
				reportsCollection:       tt.fields.reportsCollection,
			}
			got, err := m.GetUserByEmail(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("GetUserByEmail() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMongoDBRepo_GetUserByID(t *testing.T) {
	type fields struct {
		userInfoCollection      db.Collection
		passwordResetCollection db.Collection
		metadataCollection      db.Collection
		ratingsCollection       db.Collection
		faqsCollection          db.Collection
		moderateCollection      db.Collection
		reportsCollection       db.Collection
	}
	type args struct {
		id primitive.ObjectID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MongoDBRepo{
				userInfoCollection:      tt.fields.userInfoCollection,
				passwordResetCollection: tt.fields.passwordResetCollection,
				metadataCollection:      tt.fields.metadataCollection,
				ratingsCollection:       tt.fields.ratingsCollection,
				faqsCollection:          tt.fields.faqsCollection,
				moderateCollection:      tt.fields.moderateCollection,
				reportsCollection:       tt.fields.reportsCollection,
			}
			got, err := m.GetUserByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMongoDBRepo_InsertModerationData(t *testing.T) {
	type fields struct {
		userInfoCollection      db.Collection
		passwordResetCollection db.Collection
		metadataCollection      db.Collection
		ratingsCollection       db.Collection
		faqsCollection          db.Collection
		moderateCollection      db.Collection
		reportsCollection       db.Collection
	}
	type args struct {
		userID         primitive.ObjectID
		documentID     primitive.ObjectID
		approvalStatus string
		comments       string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MongoDBRepo{
				userInfoCollection:      tt.fields.userInfoCollection,
				passwordResetCollection: tt.fields.passwordResetCollection,
				metadataCollection:      tt.fields.metadataCollection,
				ratingsCollection:       tt.fields.ratingsCollection,
				faqsCollection:          tt.fields.faqsCollection,
				moderateCollection:      tt.fields.moderateCollection,
				reportsCollection:       tt.fields.reportsCollection,
			}
			if err := m.InsertModerationData(tt.args.userID, tt.args.documentID, tt.args.approvalStatus, tt.args.comments); (err != nil) != tt.wantErr {
				t.Errorf("InsertModerationData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMongoDBRepo_InsertReport(t *testing.T) {
	type fields struct {
		userInfoCollection      db.Collection
		passwordResetCollection db.Collection
		metadataCollection      db.Collection
		ratingsCollection       db.Collection
		faqsCollection          db.Collection
		moderateCollection      db.Collection
		reportsCollection       db.Collection
	}
	type args struct {
		report bson.M
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *mongo.InsertOneResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MongoDBRepo{
				userInfoCollection:      tt.fields.userInfoCollection,
				passwordResetCollection: tt.fields.passwordResetCollection,
				metadataCollection:      tt.fields.metadataCollection,
				ratingsCollection:       tt.fields.ratingsCollection,
				faqsCollection:          tt.fields.faqsCollection,
				moderateCollection:      tt.fields.moderateCollection,
				reportsCollection:       tt.fields.reportsCollection,
			}
			got, err := m.InsertReport(tt.args.report)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertReport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsertReport() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMongoDBRepo_RegisterUser(t *testing.T) {
	type fields struct {
		userInfoCollection      db.Collection
		passwordResetCollection db.Collection
		metadataCollection      db.Collection
		ratingsCollection       db.Collection
		faqsCollection          db.Collection
		moderateCollection      db.Collection
		reportsCollection       db.Collection
	}
	type args struct {
		user *models.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid user",
			args: args{user: &testUserJoe},
			fields: fields{
				userInfoCollection: &db.MongoCollectionMock{
					InsertOneFunc: func(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
						// make sure the document is a pointer to a user
						usr, ok := document.(*models.User)
						if !ok {
							return nil, mongo.ErrNilDocument
						}
						// make sure the pointer is not nil
						if usr == nil {
							return nil, mongo.ErrNilDocument
						}
						return &mongo.InsertOneResult{}, nil
					},
				},
			},
		},
		{
			name: "invalid user",
			args: args{user: nil},
			fields: fields{
				userInfoCollection: &db.MongoCollectionMock{
					InsertOneFunc: func(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
						usr, ok := document.(*models.User)
						if !ok {
							return nil, mongo.ErrNilDocument
						}
						if usr == nil {
							return nil, mongo.ErrNilDocument
						}
						return &mongo.InsertOneResult{}, nil
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MongoDBRepo{
				userInfoCollection:      tt.fields.userInfoCollection,
				passwordResetCollection: tt.fields.passwordResetCollection,
				metadataCollection:      tt.fields.metadataCollection,
				ratingsCollection:       tt.fields.ratingsCollection,
				faqsCollection:          tt.fields.faqsCollection,
				moderateCollection:      tt.fields.moderateCollection,
				reportsCollection:       tt.fields.reportsCollection,
			}
			if err := m.RegisterUser(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMongoDBRepo_SetDocumentRating(t *testing.T) {
	type fields struct {
		userInfoCollection      db.Collection
		passwordResetCollection db.Collection
		metadataCollection      db.Collection
		ratingsCollection       db.Collection
		faqsCollection          db.Collection
		moderateCollection      db.Collection
		reportsCollection       db.Collection
	}
	type args struct {
		docID  primitive.ObjectID
		rating *models.Rating
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MongoDBRepo{
				userInfoCollection:      tt.fields.userInfoCollection,
				passwordResetCollection: tt.fields.passwordResetCollection,
				metadataCollection:      tt.fields.metadataCollection,
				ratingsCollection:       tt.fields.ratingsCollection,
				faqsCollection:          tt.fields.faqsCollection,
				moderateCollection:      tt.fields.moderateCollection,
				reportsCollection:       tt.fields.reportsCollection,
			}
			if err := m.SetDocumentRating(tt.args.docID, tt.args.rating); (err != nil) != tt.wantErr {
				t.Errorf("SetDocumentRating() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMongoDBRepo_StoreResetToken(t *testing.T) {
	type fields struct {
		userInfoCollection      db.Collection
		passwordResetCollection db.Collection
		metadataCollection      db.Collection
		ratingsCollection       db.Collection
		faqsCollection          db.Collection
		moderateCollection      db.Collection
		reportsCollection       db.Collection
	}
	type args struct {
		resetEntry *models.PasswordReset
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MongoDBRepo{
				userInfoCollection:      tt.fields.userInfoCollection,
				passwordResetCollection: tt.fields.passwordResetCollection,
				metadataCollection:      tt.fields.metadataCollection,
				ratingsCollection:       tt.fields.ratingsCollection,
				faqsCollection:          tt.fields.faqsCollection,
				moderateCollection:      tt.fields.moderateCollection,
				reportsCollection:       tt.fields.reportsCollection,
			}
			if err := m.StoreResetToken(tt.args.resetEntry); (err != nil) != tt.wantErr {
				t.Errorf("StoreResetToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMongoDBRepo_UpdateDocumentsByID(t *testing.T) {
	type fields struct {
		userInfoCollection      db.Collection
		passwordResetCollection db.Collection
		metadataCollection      db.Collection
		ratingsCollection       db.Collection
		faqsCollection          db.Collection
		moderateCollection      db.Collection
		reportsCollection       db.Collection
	}
	type args struct {
		documentID primitive.ObjectID
		updateData bson.M
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MongoDBRepo{
				userInfoCollection:      tt.fields.userInfoCollection,
				passwordResetCollection: tt.fields.passwordResetCollection,
				metadataCollection:      tt.fields.metadataCollection,
				ratingsCollection:       tt.fields.ratingsCollection,
				faqsCollection:          tt.fields.faqsCollection,
				moderateCollection:      tt.fields.moderateCollection,
				reportsCollection:       tt.fields.reportsCollection,
			}
			if err := m.UpdateDocumentsByID(tt.args.documentID, tt.args.updateData); (err != nil) != tt.wantErr {
				t.Errorf("UpdateDocumentsByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMongoDBRepo_UploadDocumentMetadata(t *testing.T) {
	type fields struct {
		userInfoCollection      db.Collection
		passwordResetCollection db.Collection
		metadataCollection      db.Collection
		ratingsCollection       db.Collection
		faqsCollection          db.Collection
		moderateCollection      db.Collection
		reportsCollection       db.Collection
	}
	type args struct {
		document *models.Document
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MongoDBRepo{
				userInfoCollection:      tt.fields.userInfoCollection,
				passwordResetCollection: tt.fields.passwordResetCollection,
				metadataCollection:      tt.fields.metadataCollection,
				ratingsCollection:       tt.fields.ratingsCollection,
				faqsCollection:          tt.fields.faqsCollection,
				moderateCollection:      tt.fields.moderateCollection,
				reportsCollection:       tt.fields.reportsCollection,
			}
			if err := m.UploadDocumentMetadata(tt.args.document); (err != nil) != tt.wantErr {
				t.Errorf("UploadDocumentMetadata() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMongoDBRepo_VerifyResetToken(t *testing.T) {
	type fields struct {
		userInfoCollection      db.Collection
		passwordResetCollection db.Collection
		metadataCollection      db.Collection
		ratingsCollection       db.Collection
		faqsCollection          db.Collection
		moderateCollection      db.Collection
		reportsCollection       db.Collection
	}
	type args struct {
		userID primitive.ObjectID
		token  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MongoDBRepo{
				userInfoCollection:      tt.fields.userInfoCollection,
				passwordResetCollection: tt.fields.passwordResetCollection,
				metadataCollection:      tt.fields.metadataCollection,
				ratingsCollection:       tt.fields.ratingsCollection,
				faqsCollection:          tt.fields.faqsCollection,
				moderateCollection:      tt.fields.moderateCollection,
				reportsCollection:       tt.fields.reportsCollection,
			}
			got, err := m.VerifyResetToken(tt.args.userID, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyResetToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("VerifyResetToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMongoDBRepo(t *testing.T) {
	type args struct {
		client       *mongo.Client
		databaseName string
	}
	tests := []struct {
		name string
		args args
		want *MongoDBRepo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMongoDBRepo(tt.args.client, tt.args.databaseName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMongoDBRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}
