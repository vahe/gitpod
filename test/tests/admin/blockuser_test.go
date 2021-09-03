// Copyright (c) 2021 Gitpod GmbH. All rights reserved.
// Licensed under the GNU Affero General Public License (AGPL).
// See License-AGPL.txt in the project root for license information.

package admin

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	protocol "github.com/gitpod-io/gitpod/gitpod-protocol"
	"github.com/gitpod-io/gitpod/test/pkg/integration"
)

func TestAdminBlockUser(t *testing.T) {
	it, ctx := integration.NewTest(t, 5*time.Minute)
	defer it.Done()

	rand.Seed(time.Now().UnixNano())
	randN := rand.Intn(1000)

	adminUsername := fmt.Sprintf("admin%d", randN)
	adminUserId, err := integration.CreateUser(it, adminUsername, true)
	if err != nil {
		t.Fatalf("cannot create user: %q", err)
	}
	defer func() {
		err := integration.DeleteUser(it, adminUserId)
		if err != nil {
			t.Fatalf("error deleting user %q", err)
		}
	}()
	t.Logf("user '%s' with ID %s created", adminUsername, adminUserId)

	username := fmt.Sprintf("johndoe%d", randN)
	userId, err := integration.CreateUser(it, username, false)
	if err != nil {
		t.Fatalf("cannot create user: %q", err)
	}
	defer func() {
		err := integration.DeleteUser(it, userId)
		if err != nil {
			t.Fatalf("error deleting user %q", err)
		}
	}()
	t.Logf("user '%s' with ID %s created", username, userId)

	serverOpts := []integration.GitpodServerOpt{integration.WithGitpodUser(adminUsername)}
	server := it.API().GitpodServer(serverOpts...)
	err = server.AdminBlockUser(ctx, &protocol.AdminBlockUserRequest{UserID: userId, IsBlocked: true})
	if err != nil {
		t.Fatalf("cannot perform AdminBlockUser: %q", err)
	}

	blocked, err := integration.IsUserBlocked(it, userId)
	if err != nil {
		t.Fatalf("error checking if user is blocked: %q", err)
	}

	if !blocked {
		t.Fatalf("expected user '%s' with ID %s is blocked, but is not", username, userId)
	}
}

func TestAdminBlockUserAndStopWorkspaces(t *testing.T) {
	it, ctx := integration.NewTest(t, 5*time.Minute)
	defer it.Done()

	rand.Seed(time.Now().UnixNano())
	randN := rand.Intn(1000)

	adminUsername := fmt.Sprintf("admin%d", randN)
	adminUserId, err := integration.CreateUser(it, adminUsername, true)
	if err != nil {
		t.Fatalf("cannot create user: %q", err)
	}
	defer func() {
		err := integration.DeleteUser(it, adminUserId)
		if err != nil {
			t.Fatalf("error deleting user %q", err)
		}
	}()
	t.Logf("user '%s' with ID %s created", adminUsername, adminUserId)

	username := fmt.Sprintf("johndoe%d", randN)
	userId, err := integration.CreateUser(it, username, false)
	if err != nil {
		t.Fatalf("cannot create user: %q", err)
	}
	defer func() {
		err := integration.DeleteUser(it, userId)
		if err != nil {
			t.Fatalf("error deleting user %q", err)
		}
	}()
	t.Logf("user '%s' with ID %s created", username, userId)

	serverOpts := []integration.GitpodServerOpt{integration.WithGitpodUser(username)}
	server := it.API().GitpodServer(serverOpts...)
	// FIXME: 401 Unauthorized
	resp, err := server.CreateWorkspace(ctx, &protocol.CreateWorkspaceOptions{
		ContextURL: "github.com/gitpod-io/gitpod",
		Mode:       "force-new",
	})
	if err != nil {
		t.Fatalf("cannot start workspace: %q", err)
	}
	defer func() {
		cctx, ccancel := context.WithTimeout(context.Background(), 10*time.Second)
		err := server.StopWorkspace(cctx, resp.CreatedWorkspaceID)
		ccancel()
		if err != nil {
			t.Errorf("cannot stop workspace: %q", err)
		}
	}()
	t.Logf("created workspace: workspaceID=%s url=%s", resp.CreatedWorkspaceID, resp.WorkspaceURL)

	nfo, err := server.GetWorkspace(ctx, resp.CreatedWorkspaceID)
	if err != nil {
		t.Fatalf("cannot get workspace: %q", err)
	}
	if nfo.LatestInstance == nil {
		t.Fatal("CreateWorkspace did not start the workspace")
	}

	it.WaitForWorkspaceStart(ctx, nfo.LatestInstance.ID)

	serverOpts = []integration.GitpodServerOpt{integration.WithGitpodUser(adminUsername)}
	server = it.API().GitpodServer(serverOpts...)
	err = server.AdminBlockUser(ctx, &protocol.AdminBlockUserRequest{UserID: userId, IsBlocked: true})
	if err != nil {
		t.Fatalf("cannot perform AdminBlockUser: %q", err)
	}

	blocked, err := integration.IsUserBlocked(it, userId)
	if err != nil {
		t.Fatalf("error checking if user is blocked: %q", err)
	}

	if !blocked {
		t.Fatalf("expected user '%s' with ID %s is blocked, but is not", username, userId)
	}

	nfo, err = server.GetWorkspace(ctx, resp.CreatedWorkspaceID)
	if err != nil {
		t.Fatalf("cannot get workspace: %q", err)
	}
	if nfo.LatestInstance == nil {
		t.Fatal("CreateWorkspace did not start the workspace")
	}

	// FIXME: check for stopping status
	t.Logf("workspace status: %s", nfo.LatestInstance.Status.Phase)
}
