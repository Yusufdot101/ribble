package api

// func TestLogout(t *testing.T) {
// 	gin.SetMode(gin.TestMode)
//
// 	// create postgres test container
// 	ctx := context.Background()
// 	ctr, err := postgres.Run(
// 		ctx,
// 		"postgres:18.3-alpine",
// 		postgres.WithDatabase("ripple_user_service_test"),
// 		postgres.WithUsername("user_service"),
// 		postgres.WithPassword("verystrongpassword"),
// 		postgres.BasicWaitStrategies(),
// 	)
// 	testcontainers.CleanupContainer(t, ctr)
// 	if err != nil {
// 		log.Fatalf("failed to start postgresql: %v", err)
// 	}
//
// 	databaseURL, err := ctr.ConnectionString(ctx)
// 	if err != nil {
// 		log.Fatalf("failed to obtain connection string: %v", err)
// 	}
//
// 	repo, err := postgresql.NewAdapter(databaseURL)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	tsvc := services.NewTokenService(repo)
//
// 	// mock the google oidc
// 	provider := ports.NewMockOAuthProvider(t)
// 	svc := services.NewAuthService(repo, provider, tsvc)
//
// 	// insert user
// 	user := domain.NewUser("yusuf", "email@example.com", "google", "11")
// 	err = svc.NewUser(user)
// 	if err != nil {
// 		t.Fatalf("unexpected error inserting user: %v", err)
// 	}
//
// 	// insert token
// 	token := domain.NewToken(domain.REFRESH, domain.UUID, user.ID, "refreshToken", time.Now().Add(time.Hour))
// 	err = tsvc.Save(token)
// 	if err != nil {
// 		t.Fatalf("unexpected error inserting token: %v", err)
// 	}
//
// 	h := NewHandler(svc, tsvc, nil)
// 	resp := httptest.NewRecorder()
// 	req := httptest.NewRequest(http.MethodGet, "/auth/logout", nil)
// 	c, _ := gin.CreateTestContext(resp)
// 	c.Request = req
// 	req.AddCookie(&http.Cookie{
// 		Name:  "refreshToken",
// 		Value: token.TokenString,
// 	})
//
// 	h.logout(c)
// 	assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
// 	assert.Contains(t, resp.Body.String(), "logged out successfully")
// }
