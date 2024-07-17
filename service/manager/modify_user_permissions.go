package manager

//func (srv *akashaManagerSrvImpl) ModifyPermissions(ctx http.Context[*model.ModifyUserPermissionsRequest, *model.BaseResponse]) {
//	token, request := common.GetApiKey(ctx.HeaderParams()), ctx.Request()
//
//	// check user exist
//	user, checkErr := srv.db.GetUserByID(ctx, int64(request.UserID))
//	if checkErr != nil {
//		ctx.SetStatusCode(http.StatusNotFound)
//		ctx.SetResponse(&model.BaseResponse{
//			ErrorCode:    http.ErrorCodeResourceNotFound,
//			ErrorMessage: "user not found",
//			RequestID:    trace.GetTid(ctx),
//		})
//		return
//	}
//
//	// check providers exist
//	var permissions []string
//	for _, permission := range request.Permissions {
//		permissions = append(permissions, permission.Provider)
//	}
//
//}
