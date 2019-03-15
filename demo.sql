CREATE TABLE `user` (
	`uid` INTEGER PRIMARY KEY AUTOINCREMENT,
	`username` VARCHAR(64) NULL,
	`department` VARCHAR(64) NULL,
	`created` DATE NULL
);


INSERT INTO user(username, department, created) values('bingoohuang','bjca系统架构部','2019-03-15');
INSERT INTO user(username, department, created) values('dingoohuang','bjca电子合同部','2017-03-15');

select * from  user;
