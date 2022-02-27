package main

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"test_pet/pkg/grpc/userapi"
	"time"
)

func main() {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	client := userapi.NewUserServiceClient(conn)

	userId, err := testAddUser(client)
	if err != nil {
		log.Fatalf("failed to add user %v", err)
	}
	log.Printf("User successfully added. UserId: %d", userId)

	err = testDeleteUser(client, userId)
	if err != nil {
		log.Fatalf("failed to delete user %v", err)
	}
	log.Printf("Successfully deleted")

	err = testDeleteUser(client, userId)
	if err != missingClient {
		log.Fatalf("unexpected behaivor err: %v", err)
	}
	log.Printf("Cannot delete same client twice. success")

	if err = addFiveClients(client); err != nil {
		log.Fatalf("failed to add user %v", err)
	}

	list, err := testGetList(client)
	if err != nil {
		log.Fatalf("failed get list err: %v", err)
	}
	log.Printf("Get list result: %+v", list)
}

func testAddUser(client userapi.UserServiceClient) (int64, error) {
	addUserContext, addUserCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer addUserCancel()
	addUserInput := userapi.AddUserInput{Name: "nikitos"}
	addResult, err := client.AddUserRequest(addUserContext, &addUserInput)
	if err != nil {
		return 0, err
	}
	return addResult.GetId(), nil
}

var missingClient = errors.New("missing client")

func testDeleteUser(client userapi.UserServiceClient, userId int64) error {
	deleteUserContext, deleteUserCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer deleteUserCancel()
	deleteUserInput := userapi.DeleteUserInput{Id: userId}
	deleteResult, err := client.DeleteUserRequest(deleteUserContext, &deleteUserInput)
	if err != nil {
		return err
	}
	if deleteResult.Error != nil {
		return missingClient
	}
	return nil
}

func testGetList(client userapi.UserServiceClient) ([]*userapi.User, error) {
	getListContext, getListCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer getListCancel()
	getListInput := userapi.GetListInput{}
	getListResult, err := client.GetListRequest(getListContext, &getListInput)
	if err != nil {
		return nil, err
	}
	return getListResult.GetList(), nil
}

func addFiveClients(client userapi.UserServiceClient) error {
	_, err := testAddUser(client)
	if err != nil {
		return err
	}
	_, err = testAddUser(client)
	if err != nil {
		return err
	}
	_, err = testAddUser(client)
	if err != nil {
		return err
	}
	return nil
}
