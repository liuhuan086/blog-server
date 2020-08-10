#### 数据库语句
```mysql
create DATABASE 
if 
    not exists blog_service default character set 
    utf8mb4 default collate utf8mb4_general_ci;

create table `blog_tag` (
`id` int(10) unsigned not null auto_increment,
`name` varchar(100) default '' comment 'tag name',
`state` tinyint(3) unsigned default '1' comment '状态：0为禁用，1为启用',
`created_on` int(10) unsigned default '0' comment '创建时间',
`created_by` varchar(100) default '' comment '创建人',
`modified_on` int(10) unsigned default '0' comment '修改时间',
`modified_by` varchar(100) default '' comment '修改人',
`deleted_on` int(10) unsigned default '0' comment '删除时间',
`is_del` tinyint(3) unsigned default '0' comment '是否删除：0为未删除，1为已删除',
primary key (`id`)
)  engine=InnoDB default charset=utf8mb4 comment '标签管理';

create table `blog_article` (
`id` int(10) unsigned not null auto_increment,
`title` varchar(100) default '' comment '文章标题',
`desc` varchar(255) default '' comment '文章摘要',
`cover_image_url` varchar(255) default '' comment '封面图片地址',
`content` longtext comment '文章内容',
`created_on` int(10) unsigned default '0' comment '创建时间',
`created_by` varchar(100) default '' comment '创建人',
`modified_on` int(10) unsigned default '0' comment '修改时间',
`modified_by` varchar(100) default '' comment '修改人',
`deleted_on` int(10) unsigned default '0' comment '删除时间',
`is_del` tinyint(3) unsigned default '0' comment '是否删除：0为未删除，1为已删除',
`state` tinyint(3) unsigned default '1' comment '状态：0为禁用，1为启用',
primary key (`id`)
) engine=InnoDB default charset=utf8mb4 comment='文章管理';

create table `blog_article_tag` (
`id` int(10) unsigned not null auto_increment,
`article_id` int(11) not null comment '文章ID',
`tag_id` int(10) unsigned not null comment '标签ID',
`created_on` int(10) unsigned default '0' comment '创建时间',
`created_by` varchar(100) default '' comment '创建人',
`modified_on` int(10) unsigned default '0' comment '修改时间',
`modified_by` varchar(100) default '' comment '修改人',
`deleted_on` int(10) unsigned default '0' comment '删除时间',
`is_del` tinyint(3) unsigned default '0' comment '是否删除：0为未删除，1为已删除',
`state` tinyint(3) unsigned default '1' comment '状态：0为禁用，1为启用',
primary key (`id`)
) engine=InnoDB default charset=utf8mb4 comment='文章标签关联';
```

Verification interface
```shell script
curl -X POST http://127.0.0.1:8000/api/v1/articles -F tag_id=3 -F 'title=新增文章01-标题' -F 'desc=新增文章01-简述' \
-F 'cover_image_url=https://www.baidu.com' -F 'content=测试文章-01内容' -F 'created_by=Ryan' -F 'state=1'
```

```shell script
curl -X PUT http://127.0.0.1:8000/api/v1/articles/1 -F tag_id=3 -F 'title=测试文章-标题-更新' -F 'desc=测试文章-更新' \
 -F 'cover_image_url=https://www.baidu.com' -F 'content=测试文章-内容-更新' -F 'modified_by=Ryan' -F 'state=1'
```

```shell script
curl -X DELETE http://127.0.0.1:8000/api/v1/articles/1
```

```shell script
curl -X GET http://127.0.0.1:8000/api/v1/articles/1
```

```shell script
curl -X POST http://127.0.0.1:8000/upload/file -F file=@/Users/liuhuan/Documents/Pics/001.jpg -F type=1
```

认证表
```sql
create TABLE `blog_auth` (
`id` int(10) unsigned NOT NULL auto_increment,
`app_key` varchar(20) DEFAULT '' COMMENT 'Key',
`app_secret` varchar(50) DEFAULT '' comment 'Secret',
`created_on` int(10) unsigned default '0' comment '创建时间',
`created_by` varchar(100) default '' comment '创建人',
`modified_on` int(10) unsigned default '0' comment '修改时间',
`modified_by` varchar(100) default '' comment '修改人',
`deleted_on` int(10) unsigned default '0' comment '删除时间',
`is_del` tinyint(3) unsigned default '0' comment '是否删除：0为未删除，1为已删除',
`state` tinyint(3) unsigned default '1' comment '状态：0为禁用，1为启用',
PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT charset=utf8mb4 comment='认证管理';
```

```sql
INSERT INTO `blog_service`.`blog_auth`(`id`, `app_key`, `app_secret`, `created_on`, `created_by`, `modified_on`,
`modified_by`, `deleted_on`, `is_del`) VALUES (1, 'Ryan', 'go-programming-tour', 0, 'Ryan',0,'',0,0);
```

```shell script
curl -X GET 'http://127.0.0.1:8000/auth?app_key=xxx&app_secret=yyy'
```