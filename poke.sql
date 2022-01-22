/*
SQLyog Ultimate
MySQL - 10.4.20-MariaDB : Database - poke
*********************************************************************
*/

/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
CREATE DATABASE /*!32312 IF NOT EXISTS*/`poke` /*!40100 DEFAULT CHARACTER SET utf8mb4 */;

USE `poke`;

/*Table structure for table `game` */

DROP TABLE IF EXISTS `game`;

CREATE TABLE `game` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `game_name` varchar(255) DEFAULT 'default',
  `start_time` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

/*Data for the table `game` */

insert  into `game`(`id`,`game_name`,`start_time`) values (1,'default','2021-08-08 11:38:58');

/*Table structure for table `game_match` */

DROP TABLE IF EXISTS `game_match`;

CREATE TABLE `game_match` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `game_id` int(11) DEFAULT NULL,
  `small_bind_position` int(11) DEFAULT NULL COMMENT '小盲位',
  `big_bind_position` int(11) DEFAULT NULL COMMENT '大盲位',
  `pot_all` int(11) DEFAULT NULL COMMENT '底池',
  `pot1st` int(11) DEFAULT NULL COMMENT '第一轮底池',
  `pot2nd` int(11) DEFAULT NULL COMMENT '第二轮底池',
  `pot3rd` int(11) DEFAULT NULL COMMENT '第三轮底池',
  `pot4th` int(11) DEFAULT NULL COMMENT '第四轮地址',
  `game_status` enum('INIT','LICENSING','ROUND1','ROUND2','ROUND3','ROUND4','END') DEFAULT NULL,
  `create_time` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8;

/*Data for the table `game_match` */

insert  into `game_match`(`id`,`game_id`,`small_bind_position`,`big_bind_position`,`pot_all`,`pot1st`,`pot2nd`,`pot3rd`,`pot4th`,`game_status`,`create_time`) values (1,1,1,2,0,0,0,0,0,'ROUND1','2022-01-09 18:20:31');
insert  into `game_match`(`id`,`game_id`,`small_bind_position`,`big_bind_position`,`pot_all`,`pot1st`,`pot2nd`,`pot3rd`,`pot4th`,`game_status`,`create_time`) values (2,2,1,2,0,0,0,0,0,'ROUND1','2022-01-09 18:22:04');
insert  into `game_match`(`id`,`game_id`,`small_bind_position`,`big_bind_position`,`pot_all`,`pot1st`,`pot2nd`,`pot3rd`,`pot4th`,`game_status`,`create_time`) values (3,3,1,2,20,20,0,0,0,'END','2022-01-09 22:50:48');
insert  into `game_match`(`id`,`game_id`,`small_bind_position`,`big_bind_position`,`pot_all`,`pot1st`,`pot2nd`,`pot3rd`,`pot4th`,`game_status`,`create_time`) values (4,3,2,1,0,0,0,0,0,'ROUND1','2022-01-09 23:28:34');
insert  into `game_match`(`id`,`game_id`,`small_bind_position`,`big_bind_position`,`pot_all`,`pot1st`,`pot2nd`,`pot3rd`,`pot4th`,`game_status`,`create_time`) values (5,4,1,2,20,20,0,0,0,'END','2022-01-11 22:31:09');
insert  into `game_match`(`id`,`game_id`,`small_bind_position`,`big_bind_position`,`pot_all`,`pot1st`,`pot2nd`,`pot3rd`,`pot4th`,`game_status`,`create_time`) values (6,4,2,1,20,20,0,0,0,'END','2022-01-11 22:33:06');
insert  into `game_match`(`id`,`game_id`,`small_bind_position`,`big_bind_position`,`pot_all`,`pot1st`,`pot2nd`,`pot3rd`,`pot4th`,`game_status`,`create_time`) values (7,4,1,2,0,0,0,0,0,'ROUND1','2022-01-11 22:33:32');
insert  into `game_match`(`id`,`game_id`,`small_bind_position`,`big_bind_position`,`pot_all`,`pot1st`,`pot2nd`,`pot3rd`,`pot4th`,`game_status`,`create_time`) values (8,5,1,2,110,110,0,0,0,'END','2022-01-11 22:36:24');
insert  into `game_match`(`id`,`game_id`,`small_bind_position`,`big_bind_position`,`pot_all`,`pot1st`,`pot2nd`,`pot3rd`,`pot4th`,`game_status`,`create_time`) values (9,5,2,1,20,20,0,0,0,'END','2022-01-11 22:43:10');
insert  into `game_match`(`id`,`game_id`,`small_bind_position`,`big_bind_position`,`pot_all`,`pot1st`,`pot2nd`,`pot3rd`,`pot4th`,`game_status`,`create_time`) values (10,5,1,2,20,20,0,0,0,'END','2022-01-11 22:52:16');
insert  into `game_match`(`id`,`game_id`,`small_bind_position`,`big_bind_position`,`pot_all`,`pot1st`,`pot2nd`,`pot3rd`,`pot4th`,`game_status`,`create_time`) values (11,5,2,1,20,20,0,0,0,'END','2022-01-11 22:52:36');
insert  into `game_match`(`id`,`game_id`,`small_bind_position`,`big_bind_position`,`pot_all`,`pot1st`,`pot2nd`,`pot3rd`,`pot4th`,`game_status`,`create_time`) values (12,5,1,2,0,0,0,0,0,'ROUND1','2022-01-11 22:55:01');
insert  into `game_match`(`id`,`game_id`,`small_bind_position`,`big_bind_position`,`pot_all`,`pot1st`,`pot2nd`,`pot3rd`,`pot4th`,`game_status`,`create_time`) values (13,6,1,2,0,0,0,0,0,'ROUND1','2022-01-11 22:56:00');
insert  into `game_match`(`id`,`game_id`,`small_bind_position`,`big_bind_position`,`pot_all`,`pot1st`,`pot2nd`,`pot3rd`,`pot4th`,`game_status`,`create_time`) values (14,7,1,2,20,20,0,0,0,'END','2022-01-12 19:59:25');
insert  into `game_match`(`id`,`game_id`,`small_bind_position`,`big_bind_position`,`pot_all`,`pot1st`,`pot2nd`,`pot3rd`,`pot4th`,`game_status`,`create_time`) values (15,7,2,1,20,20,0,0,0,'END','2022-01-12 20:04:00');
insert  into `game_match`(`id`,`game_id`,`small_bind_position`,`big_bind_position`,`pot_all`,`pot1st`,`pot2nd`,`pot3rd`,`pot4th`,`game_status`,`create_time`) values (16,7,1,2,0,0,0,0,0,'ROUND1','2022-01-12 20:04:46');
insert  into `game_match`(`id`,`game_id`,`small_bind_position`,`big_bind_position`,`pot_all`,`pot1st`,`pot2nd`,`pot3rd`,`pot4th`,`game_status`,`create_time`) values (17,8,1,2,20,20,0,0,0,'END','2022-01-12 20:06:09');
insert  into `game_match`(`id`,`game_id`,`small_bind_position`,`big_bind_position`,`pot_all`,`pot1st`,`pot2nd`,`pot3rd`,`pot4th`,`game_status`,`create_time`) values (18,8,2,1,20,20,0,0,0,'END','2022-01-12 20:07:29');
insert  into `game_match`(`id`,`game_id`,`small_bind_position`,`big_bind_position`,`pot_all`,`pot1st`,`pot2nd`,`pot3rd`,`pot4th`,`game_status`,`create_time`) values (19,8,1,2,0,0,0,0,0,'ROUND1','2022-01-12 20:08:02');

/*Table structure for table `game_match_log` */

DROP TABLE IF EXISTS `game_match_log`;

CREATE TABLE `game_match_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `game_match_id` int(11) DEFAULT NULL,
  `user_id` int(11) DEFAULT NULL,
  `operation` varchar(255) DEFAULT NULL,
  `point_number` int(11) DEFAULT NULL,
  `add_time` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=132 DEFAULT CHARSET=utf8mb4;

/*Data for the table `game_match_log` */

insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (1,1,1,'raise',5,'2022-01-09 18:20:31');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (2,1,0,'raise',10,'2022-01-09 18:20:31');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (3,2,1,'raise',5,'2022-01-09 18:22:05');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (4,2,0,'raise',10,'2022-01-09 18:22:05');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (5,3,2,'raise',5,'2022-01-09 22:50:48');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (6,3,0,'raise',10,'2022-01-09 22:50:48');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (7,3,2,'call',5,'2022-01-09 22:50:53');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (8,3,1,'call',0,'2022-01-09 22:50:58');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (9,3,2,'check',0,'2022-01-09 22:51:34');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (10,3,1,'check',0,'2022-01-09 22:51:37');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (11,3,2,'check',0,'2022-01-09 22:51:48');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (12,3,1,'check',0,'2022-01-09 22:51:51');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (13,3,2,'check',0,'2022-01-09 22:52:10');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (14,3,1,'check',0,'2022-01-09 22:52:12');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (15,4,1,'raise',5,'2022-01-09 23:28:35');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (16,4,0,'raise',10,'2022-01-09 23:28:35');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (17,5,1,'raise',5,'2022-01-11 22:31:09');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (18,5,0,'raise',10,'2022-01-11 22:31:09');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (19,5,1,'raise',0,'2022-01-11 22:31:30');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (20,5,2,'call',0,'2022-01-11 22:31:43');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (21,5,1,'call',5,'2022-01-11 22:31:47');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (22,5,1,'check',0,'2022-01-11 22:31:54');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (23,5,2,'check',0,'2022-01-11 22:31:57');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (24,5,1,'check',0,'2022-01-11 22:31:59');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (25,5,2,'check',0,'2022-01-11 22:32:01');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (26,5,1,'check',0,'2022-01-11 22:32:03');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (27,5,2,'check',0,'2022-01-11 22:32:04');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (28,6,2,'raise',5,'2022-01-11 22:33:06');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (29,6,0,'raise',10,'2022-01-11 22:33:07');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (30,6,2,'call',5,'2022-01-11 22:33:08');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (31,6,1,'call',0,'2022-01-11 22:33:09');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (32,6,2,'call',0,'2022-01-11 22:33:11');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (33,6,2,'check',0,'2022-01-11 22:33:11');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (34,6,1,'check',0,'2022-01-11 22:33:13');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (35,6,2,'check',0,'2022-01-11 22:33:14');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (36,6,1,'check',0,'2022-01-11 22:33:16');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (37,7,1,'raise',5,'2022-01-11 22:33:32');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (38,7,0,'raise',10,'2022-01-11 22:33:32');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (39,8,2,'raise',5,'2022-01-11 22:36:25');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (40,8,0,'raise',10,'2022-01-11 22:36:25');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (41,8,2,'raise',50,'2022-01-11 22:41:39');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (42,8,1,'call',45,'2022-01-11 22:41:51');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (43,8,2,'check',0,'2022-01-11 22:42:01');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (44,8,1,'check',0,'2022-01-11 22:42:04');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (45,8,2,'check',0,'2022-01-11 22:42:06');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (46,8,1,'check',0,'2022-01-11 22:42:08');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (47,8,2,'raise',50,'2022-01-11 22:42:15');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (48,8,1,'check',0,'2022-01-11 22:42:22');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (49,8,2,'check',0,'2022-01-11 22:42:25');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (50,9,1,'raise',5,'2022-01-11 22:43:10');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (51,9,0,'raise',10,'2022-01-11 22:43:10');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (52,9,1,'call',5,'2022-01-11 22:45:55');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (53,9,2,'call',0,'2022-01-11 22:46:07');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (54,9,1,'check',0,'2022-01-11 22:46:19');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (55,9,2,'check',0,'2022-01-11 22:46:21');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (56,9,1,'fold',0,'2022-01-11 22:46:25');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (57,9,1,'fold',0,'2022-01-11 22:46:33');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (58,9,1,'fold',0,'2022-01-11 22:46:34');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (59,9,1,'fold',0,'2022-01-11 22:46:34');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (60,9,1,'fold',0,'2022-01-11 22:46:35');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (61,9,1,'fold',0,'2022-01-11 22:46:35');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (62,9,1,'fold',0,'2022-01-11 22:46:35');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (63,9,1,'fold',0,'2022-01-11 22:46:36');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (64,9,1,'fold',0,'2022-01-11 22:46:36');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (65,9,1,'fold',0,'2022-01-11 22:46:36');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (66,9,1,'fold',0,'2022-01-11 22:46:37');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (67,9,1,'fold',0,'2022-01-11 22:47:17');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (68,10,2,'raise',5,'2022-01-11 22:52:16');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (69,10,0,'raise',10,'2022-01-11 22:52:16');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (70,10,2,'call',5,'2022-01-11 22:52:20');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (71,10,1,'call',0,'2022-01-11 22:52:24');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (72,10,2,'check',0,'2022-01-11 22:52:26');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (73,10,1,'check',0,'2022-01-11 22:52:27');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (74,10,2,'check',0,'2022-01-11 22:52:28');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (75,10,1,'check',0,'2022-01-11 22:52:30');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (76,10,2,'check',0,'2022-01-11 22:52:31');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (77,10,1,'check',0,'2022-01-11 22:52:32');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (78,11,1,'raise',5,'2022-01-11 22:52:36');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (79,11,0,'raise',10,'2022-01-11 22:52:36');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (80,11,1,'call',5,'2022-01-11 22:52:39');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (81,11,2,'call',0,'2022-01-11 22:52:42');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (82,11,1,'fold',0,'2022-01-11 22:52:45');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (83,12,2,'raise',5,'2022-01-11 22:55:01');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (84,12,0,'raise',10,'2022-01-11 22:55:01');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (85,13,1,'raise',5,'2022-01-11 22:56:00');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (86,13,0,'raise',10,'2022-01-11 22:56:00');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (87,13,1,'call',5,'2022-01-11 22:56:05');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (88,14,2,'raise',5,'2022-01-12 19:59:25');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (89,14,0,'raise',10,'2022-01-12 19:59:25');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (90,14,2,'call',5,'2022-01-12 19:59:58');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (91,14,1,'call',0,'2022-01-12 20:00:15');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (92,14,2,'check',0,'2022-01-12 20:00:26');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (93,14,1,'check',0,'2022-01-12 20:00:38');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (94,14,2,'check',0,'2022-01-12 20:00:41');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (95,14,1,'check',0,'2022-01-12 20:00:43');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (96,14,2,'check',0,'2022-01-12 20:00:47');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (97,14,1,'check',0,'2022-01-12 20:00:49');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (98,15,1,'raise',5,'2022-01-12 20:04:00');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (99,15,0,'raise',10,'2022-01-12 20:04:00');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (100,15,1,'call',5,'2022-01-12 20:04:06');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (101,15,2,'call',0,'2022-01-12 20:04:07');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (102,15,1,'check',0,'2022-01-12 20:04:09');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (103,15,2,'check',0,'2022-01-12 20:04:10');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (104,15,1,'check',0,'2022-01-12 20:04:11');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (105,15,2,'check',0,'2022-01-12 20:04:12');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (106,15,1,'check',0,'2022-01-12 20:04:14');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (107,15,2,'check',0,'2022-01-12 20:04:15');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (108,16,2,'raise',5,'2022-01-12 20:04:46');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (109,16,0,'raise',10,'2022-01-12 20:04:46');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (110,17,2,'raise',5,'2022-01-12 20:06:09');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (111,17,0,'raise',10,'2022-01-12 20:06:09');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (112,17,2,'call',5,'2022-01-12 20:06:24');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (113,17,1,'call',0,'2022-01-12 20:06:27');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (114,17,2,'check',0,'2022-01-12 20:06:44');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (115,17,1,'check',0,'2022-01-12 20:06:48');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (116,17,2,'check',0,'2022-01-12 20:06:54');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (117,17,1,'check',0,'2022-01-12 20:06:56');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (118,17,2,'check',0,'2022-01-12 20:07:12');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (119,17,1,'check',0,'2022-01-12 20:07:14');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (120,18,1,'raise',5,'2022-01-12 20:07:29');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (121,18,0,'raise',10,'2022-01-12 20:07:29');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (122,18,1,'call',5,'2022-01-12 20:07:31');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (123,18,2,'call',0,'2022-01-12 20:07:33');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (124,18,1,'check',0,'2022-01-12 20:07:34');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (125,18,2,'check',0,'2022-01-12 20:07:35');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (126,18,1,'check',0,'2022-01-12 20:07:36');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (127,18,2,'check',0,'2022-01-12 20:07:37');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (128,18,1,'check',0,'2022-01-12 20:07:38');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (129,18,2,'check',0,'2022-01-12 20:07:39');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (130,19,2,'raise',5,'2022-01-12 20:08:02');
insert  into `game_match_log`(`id`,`game_match_id`,`user_id`,`operation`,`point_number`,`add_time`) values (131,19,0,'raise',10,'2022-01-12 20:08:02');

/*Table structure for table `game_user` */

DROP TABLE IF EXISTS `game_user`;

CREATE TABLE `game_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `game_id` int(11) DEFAULT NULL,
  `add_time` datetime DEFAULT current_timestamp(),
  `position` int(11) DEFAULT NULL,
  `online` tinyint(4) DEFAULT 1 COMMENT '1为在线',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id` (`user_id`,`game_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

/*Data for the table `game_user` */

insert  into `game_user`(`id`,`user_id`,`game_id`,`add_time`,`position`,`online`) values (1,2,8,'2022-01-12 20:06:01',1,1);
insert  into `game_user`(`id`,`user_id`,`game_id`,`add_time`,`position`,`online`) values (2,1,8,'2022-01-12 20:06:05',2,1);

/*Table structure for table `room` */

DROP TABLE IF EXISTS `room`;

CREATE TABLE `room` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `game_id` int(11) DEFAULT NULL,
  `create_user_id` int(11) DEFAULT NULL,
  `room_name` varchar(255) DEFAULT NULL,
  `room_password` varchar(255) DEFAULT NULL,
  `limit_member` int(11) DEFAULT 8,
  `init_base_point` int(11) DEFAULT 1,
  `uppoint_after_game_match_number` int(11) DEFAULT 5,
  `online` tinyint(2) DEFAULT 1,
  `card_type` enum('long','short') DEFAULT 'short',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;

/*Data for the table `room` */

insert  into `room`(`id`,`game_id`,`create_user_id`,`room_name`,`room_password`,`limit_member`,`init_base_point`,`uppoint_after_game_match_number`,`online`,`card_type`) values (1,NULL,2,'','',8,1,5,1,'short');
insert  into `room`(`id`,`game_id`,`create_user_id`,`room_name`,`room_password`,`limit_member`,`init_base_point`,`uppoint_after_game_match_number`,`online`,`card_type`) values (2,NULL,2,'','',8,1,5,1,'short');
insert  into `room`(`id`,`game_id`,`create_user_id`,`room_name`,`room_password`,`limit_member`,`init_base_point`,`uppoint_after_game_match_number`,`online`,`card_type`) values (3,NULL,2,'','',8,1,5,1,'short');
insert  into `room`(`id`,`game_id`,`create_user_id`,`room_name`,`room_password`,`limit_member`,`init_base_point`,`uppoint_after_game_match_number`,`online`,`card_type`) values (4,NULL,1,'','',8,1,5,1,'short');
insert  into `room`(`id`,`game_id`,`create_user_id`,`room_name`,`room_password`,`limit_member`,`init_base_point`,`uppoint_after_game_match_number`,`online`,`card_type`) values (5,NULL,2,'','',8,1,5,1,'short');
insert  into `room`(`id`,`game_id`,`create_user_id`,`room_name`,`room_password`,`limit_member`,`init_base_point`,`uppoint_after_game_match_number`,`online`,`card_type`) values (6,NULL,1,'','',8,1,5,1,'short');
insert  into `room`(`id`,`game_id`,`create_user_id`,`room_name`,`room_password`,`limit_member`,`init_base_point`,`uppoint_after_game_match_number`,`online`,`card_type`) values (7,NULL,2,'','',8,1,5,1,'short');
insert  into `room`(`id`,`game_id`,`create_user_id`,`room_name`,`room_password`,`limit_member`,`init_base_point`,`uppoint_after_game_match_number`,`online`,`card_type`) values (8,NULL,2,'','',8,1,5,1,'short');

/*Table structure for table `user` */

DROP TABLE IF EXISTS `user`;

CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `password` varchar(255) DEFAULT NULL,
  `point` int(11) DEFAULT NULL,
  `create_time` datetime DEFAULT current_timestamp(),
  `update_time` datetime DEFAULT current_timestamp(),
  `last_login_time` datetime DEFAULT current_timestamp(),
  `experience` bigint(20) DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4;

/*Data for the table `user` */

insert  into `user`(`id`,`name`,`password`,`point`,`create_time`,`update_time`,`last_login_time`,`experience`) values (1,'suchot','83ea91541404262bcc99877e06789115',7089,'2021-08-07 17:31:52','2021-08-07 17:31:52','2021-08-28 20:46:54',0);
insert  into `user`(`id`,`name`,`password`,`point`,`create_time`,`update_time`,`last_login_time`,`experience`) values (2,'wuke','83ea91541404262bcc99877e06789115',910,'2021-08-07 17:31:58','2021-08-07 17:31:58','2021-08-28 20:46:54',0);
insert  into `user`(`id`,`name`,`password`,`point`,`create_time`,`update_time`,`last_login_time`,`experience`) values (3,'cjx','83ea91541404262bcc99877e06789115',965,'2021-08-07 17:32:04','2021-08-07 17:32:04','2021-08-28 20:46:54',0);
insert  into `user`(`id`,`name`,`password`,`point`,`create_time`,`update_time`,`last_login_time`,`experience`) values (4,'djs','83ea91541404262bcc99877e06789115',1000,'2021-08-07 17:32:10','2021-08-07 17:32:10','2021-08-28 20:46:54',0);
insert  into `user`(`id`,`name`,`password`,`point`,`create_time`,`update_time`,`last_login_time`,`experience`) values (5,'muge','83ea91541404262bcc99877e06789115',1000,'2021-08-07 17:32:17','2021-08-07 17:32:17','2021-08-28 20:46:54',0);
insert  into `user`(`id`,`name`,`password`,`point`,`create_time`,`update_time`,`last_login_time`,`experience`) values (6,'jt','83ea91541404262bcc99877e06789115',1000,'2021-08-07 17:32:23','2021-08-07 17:32:23','2021-08-28 20:46:54',0);
insert  into `user`(`id`,`name`,`password`,`point`,`create_time`,`update_time`,`last_login_time`,`experience`) values (7,'laowang','83ea91541404262bcc99877e06789115',1500,'2021-08-07 17:32:30','2021-08-07 17:32:30','2021-08-28 20:46:54',0);
insert  into `user`(`id`,`name`,`password`,`point`,`create_time`,`update_time`,`last_login_time`,`experience`) values (8,'ashu','83ea91541404262bcc99877e06789115',0,'2021-08-07 17:32:33','2021-08-07 17:32:33','2021-08-28 20:46:54',0);
insert  into `user`(`id`,`name`,`password`,`point`,`create_time`,`update_time`,`last_login_time`,`experience`) values (9,'suchot2','83ea91541404262bcc99877e06789115',0,'2021-10-26 22:13:51','2021-10-26 22:13:51','2021-10-26 22:13:51',0);

/*Table structure for table `user_game_log` */

DROP TABLE IF EXISTS `user_game_log`;

CREATE TABLE `user_game_log` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `game_match_id` int(11) DEFAULT NULL,
  `end_point` int(11) DEFAULT 0,
  `get_experience` int(11) DEFAULT NULL,
  `end_time` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*Data for the table `user_game_log` */

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
