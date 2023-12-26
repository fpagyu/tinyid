CREATE TABLE `id_info` (
	`id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
	`biz` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '业务类型',
	`max_id` bigint(20) NOT NULL COMMENT '当前号段最大值',
	`step` int(11) NOT NULL COMMENT '号段长度',
	`version` int(11) NOT NULL COMMENT '乐观锁',
	`pid` int NOT NULL COMMENT '分片id',
	`pid_bits` tinyint NOT NULL COMMENT '分片id的二进制位数',
	PRIMARY KEY (`id`),
	Unique KEY `idx_refbizuniq`(`biz`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci COMMENT='分布式id生成器号段信息';

