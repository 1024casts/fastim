
# pro tip: user_id 表示记录的拥有者

# 消息内容表
create table `im_msg` (
    `id` bigint(20) unsigned not null auto_increment,
    `user_id` bigint(20) unsigned not null comment '用户id',
    `msg_type` tinyint(4) unsigned not null default 0 comment '消息类型',
    `receive_type` tinyint(4) unsigned not null default 0 comment '接收类型 0:双方接收 1:对方接收，2:自己接收',
    `content` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' comment '消息内容',
    `extra` varchar(200) default '' comment '扩展字段，json格式',
    `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    primary key (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='消息内容表';

# 消息会话表-和多个人的聊天列表
create table `im_chat` (
    `id` bigint(20) unsigned not null auto_increment,
    `sender_uid` bigint(20) unsigned not null comment '发送者uid',
    `receive_uid` bigint(20) unsigned not null comment '接收者uid',
    `last_msg_id` bigint(20) unsigned not null comment '最后一条的消息id',
    `msg_num` int(10) unsigned not null default 0 comment '未读消息数',
    `is_delete` tinyint(4) unsigned NOT NULL DEFAULT '0' comment '是否删除 0:否 1:是',
    `extra` varchar(200) default '' comment '扩展字段，json格式',
    `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    primary key (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='消息会话表';

# 消息用户会话索引表
# 作用1：发送消息时，根据此表查询会话表
create table `im_user_chat` (
   `id` bigint(20) unsigned not null auto_increment,
   `sender_uid` bigint(20) unsigned not null comment '发送者uid',
   `receive_uid` bigint(20) unsigned not null comment '接收者uid',
   `last_msg_id` bigint(20) unsigned not null comment '最后一条的消息id',
   `msg_num` int(10) unsigned not null default 0 comment '未读消息数',
   `extra` varchar(200) default '' comment '扩展字段，json格式',
   `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
   primary key (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='用户会话索引表';