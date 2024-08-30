package repository_test

//func Test_userRepository_Add(t *testing.T) {
//	db, mock, err := sqlxmock.Newx()
//	if err != nil {
//		panic(err)
//	}
//	defer db.Close()
//
//	type args struct {
//		ctx         context.Context
//		credentials models.Credentials
//	}
//	type want struct {
//		user *models.User
//		err  error
//	}
//	tests := []struct {
//		name string
//		mock func() UserRepository
//		args args
//		want want
//	}{
//		{
//			name: "Something wrong",
//			args: args{
//				ctx: context.Background(),
//				credentials: models.Credentials{
//					Username: []byte("username"),
//					Password: []byte("password"),
//				},
//			},
//			want: want{
//				user: nil,
//				err:  errors.New("something wrong"),
//			},
//			mock: func() UserRepository {
//				repo := NewUserRepository(db)
//
//				mock.
//					ExpectQuery("^INSERT INTO users \\(username, password\\) VALUES \\(\\$1, \\$2\\) RETURNING id$").
//					WillReturnError(errors.New("something wrong"))
//
//				return repo
//			},
//		},
//		{
//			name: "Login is busy",
//			args: args{
//				ctx: context.Background(),
//				credentials: models.Credentials{
//					Username: []byte("username"),
//					Password: []byte("password"),
//				},
//			},
//			want: want{
//				user: nil,
//				err:  ErrorBusyLogin,
//			},
//			mock: func() UserRepository {
//				repo := NewUserRepository(db)
//
//				mock.
//					ExpectQuery("^INSERT INTO users \\(username, password\\) VALUES \\(\\$1, \\$2\\) RETURNING id$").
//					WillReturnError(&pgconn.PgError{Code: "23505"})
//
//				return repo
//			},
//		},
//		{
//			name: "Successfully case",
//			args: args{
//				ctx: context.Background(),
//				credentials: models.Credentials{
//					Username: []byte("username"),
//					Password: []byte("password"),
//				},
//			},
//			want: want{
//				user: &models.User{
//					ID:       1,
//					Username: []byte("username"),
//					Password: []byte("password"),
//				},
//				err: nil,
//			},
//			mock: func() UserRepository {
//				repo := NewUserRepository(db)
//
//				rows := mock.NewRows([]string{"id"}).AddRow(1)
//
//				mock.
//					ExpectQuery("^INSERT INTO users \\(username, password\\) VALUES \\(\\$1, \\$2\\) RETURNING id$").
//					WillReturnRows(rows)
//
//				return repo
//			},
//		},
//	}
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			repo := test.mock()
//			user, err := repo.AddUser(test.args.ctx, test.args.credentials)
//			require.Equal(t, test.want.err, err)
//			require.Equal(t, test.want.user, user)
//		})
//	}
//}
//
//func Test_userRepository_FindByID(t *testing.T) {
//	db, mock, err := sqlxmock.Newx()
//	if err != nil {
//		panic(err)
//	}
//	defer db.Close()
//
//	type args struct {
//		ctx context.Context
//		id  models.UserID
//	}
//	type want struct {
//		user *models.User
//		err  error
//	}
//	tests := []struct {
//		name string
//		mock func() UserRepository
//		args args
//		want want
//	}{
//		{
//			name: "Something wrong",
//			mock: func() UserRepository {
//				repo := NewUserRepository(db)
//
//				mock.
//					ExpectQuery("^SELECT \\* FROM users WHERE id = \\$1 LIMIT 1$").
//					WillReturnError(errors.New("something wrong"))
//
//				return repo
//			},
//			args: args{
//				ctx: context.Background(),
//				id:  1,
//			},
//			want: want{
//				user: nil,
//				err:  errors.New("something wrong"),
//			},
//		},
//		{
//			name: "User not found",
//			mock: func() UserRepository {
//				repo := NewUserRepository(db)
//
//				mock.
//					ExpectQuery("^SELECT \\* FROM users WHERE id = \\$1 LIMIT 1$").
//					WillReturnError(sql.ErrNoRows)
//
//				return repo
//			},
//			args: args{
//				ctx: context.Background(),
//				id:  1,
//			},
//			want: want{
//				user: nil,
//				err:  ErrorEntityNotFound,
//			},
//		},
//		{
//			name: "Successfully case",
//			mock: func() UserRepository {
//				repo := NewUserRepository(db)
//
//				rows := sqlxmock.NewRows([]string{"id", "username", "password"}).
//					AddRow("1", "username", "password")
//				mock.
//					ExpectQuery("^SELECT \\* FROM users WHERE id = \\$1 LIMIT 1$").
//					WillReturnRows(rows)
//
//				return repo
//			},
//			args: args{
//				ctx: context.Background(),
//				id:  1,
//			},
//			want: want{
//				user: &models.User{
//					ID:       1,
//					Username: []byte("username"),
//					Password: []byte("password"),
//				},
//				err: nil,
//			},
//		},
//	}
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			repo := test.mock()
//			user, err := repo.FindUserByID(test.args.ctx, test.args.id)
//			assert.Equal(t, test.want.err, err)
//			assert.Equal(t, test.want.user, user)
//		})
//	}
//}
//
//func Test_userRepository_FindByUsername(t *testing.T) {
//	db, mock, err := sqlxmock.Newx()
//	if err != nil {
//		panic(err)
//	}
//	defer db.Close()
//
//	type args struct {
//		ctx      context.Context
//		username []byte
//	}
//	type want struct {
//		user *models.User
//		err  error
//	}
//	tests := []struct {
//		name string
//		mock func() UserRepository
//		args args
//		want want
//	}{
//		{
//			name: "Something wrong",
//			mock: func() UserRepository {
//				repo := NewUserRepository(db)
//
//				mock.
//					ExpectQuery("^SELECT \\* FROM users WHERE username = \\$1 LIMIT 1$").
//					WillReturnError(errors.New("something wrong"))
//
//				return repo
//			},
//			args: args{
//				ctx:      context.Background(),
//				username: []byte("username"),
//			},
//			want: want{
//				user: nil,
//				err:  errors.New("something wrong"),
//			},
//		},
//		{
//			name: "User not found",
//			mock: func() UserRepository {
//				repo := NewUserRepository(db)
//
//				mock.
//					ExpectQuery("^SELECT \\* FROM users WHERE username = \\$1 LIMIT 1$").
//					WillReturnError(sql.ErrNoRows)
//
//				return repo
//			},
//			args: args{
//				ctx:      context.Background(),
//				username: []byte("username"),
//			},
//			want: want{
//				user: nil,
//				err:  ErrorEntityNotFound,
//			},
//		},
//		{
//			name: "Successfully case",
//			mock: func() UserRepository {
//				repo := NewUserRepository(db)
//
//				rows := sqlxmock.NewRows([]string{"id", "username", "password"}).
//					AddRow("1", "username", "password")
//				mock.
//					ExpectQuery("^SELECT \\* FROM users WHERE username = \\$1 LIMIT 1$").
//					WillReturnRows(rows)
//
//				return repo
//			},
//			args: args{
//				ctx:      context.Background(),
//				username: []byte("username"),
//			},
//			want: want{
//				user: &models.User{
//					ID:       1,
//					Username: []byte("username"),
//					Password: []byte("password"),
//				},
//				err: nil,
//			},
//		},
//	}
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			repo := test.mock()
//			user, err := repo.FindUserByUsername(test.args.ctx, test.args.username)
//			assert.Equal(t, test.want.err, err)
//			assert.Equal(t, test.want.user, user)
//		})
//	}
//}
//
//func TestNewUserRepository(t *testing.T) {
//	db, _, err := sqlxmock.Newx()
//	if err != nil {
//		panic(err)
//	}
//
//	repo := NewUserRepository(db)
//	assert.Implements(t, (*UserRepository)(nil), repo)
//}
