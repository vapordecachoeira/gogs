// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package repo

import (
	"io"
	"path"

	"github.com/gogits/git-module"

	"github.com/gogits/gogs/pkg/tool"
	"github.com/gogits/gogs/pkg/context"
	"github.com/gogits/gogs/pkg/setting"
)

func ServeData(ctx *context.Context, name string, reader io.Reader) error {
	buf := make([]byte, 1024)
	n, _ := reader.Read(buf)
	if n > 0 {
		buf = buf[:n]
	}

	if !tool.IsTextFile(buf) {
		if !tool.IsImageFile(buf) {
			ctx.Resp.Header().Set("Content-Disposition", "attachment; filename=\""+name+"\"")
			ctx.Resp.Header().Set("Content-Transfer-Encoding", "binary")
		}
	} else if !setting.Repository.EnableRawFileRenderMode || !ctx.QueryBool("render") {
		ctx.Resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}
	ctx.Resp.Write(buf)
	_, err := io.Copy(ctx.Resp, reader)
	return err
}

func ServeBlob(ctx *context.Context, blob *git.Blob) error {
	dataRc, err := blob.Data()
	if err != nil {
		return err
	}

	return ServeData(ctx, path.Base(ctx.Repo.TreePath), dataRc)
}

func SingleDownload(ctx *context.Context) {
	blob, err := ctx.Repo.Commit.GetBlobByPath(ctx.Repo.TreePath)
	if err != nil {
		if git.IsErrNotExist(err) {
			ctx.Handle(404, "GetBlobByPath", nil)
		} else {
			ctx.Handle(500, "GetBlobByPath", err)
		}
		return
	}
	if err = ServeBlob(ctx, blob); err != nil {
		ctx.Handle(500, "ServeBlob", err)
	}
}
