package cmd

import (
	_ "github.com/shichen437/live-dog/internal/app/admin/logic"
	_ "github.com/shichen437/live-dog/internal/app/common/logic"
	_ "github.com/shichen437/live-dog/internal/app/live/logic"
	_ "github.com/shichen437/live-dog/internal/app/monitor/logic"

	_ "github.com/shichen437/live-dog/internal/pkg/lives/bilibili"
	_ "github.com/shichen437/live-dog/internal/pkg/lives/douyin"

	_ "github.com/shichen437/live-dog/internal/pkg/message_push/email"
	_ "github.com/shichen437/live-dog/internal/pkg/message_push/gotify"

	_ "github.com/shichen437/live-dog/internal/pkg/media_parser/bilibili"
	_ "github.com/shichen437/live-dog/internal/pkg/media_parser/douyin"

	_ "github.com/shichen437/live-dog/internal/pkg/download/bilibili"
	_ "github.com/shichen437/live-dog/internal/pkg/download/douyin"
)
