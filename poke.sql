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
) ENGINE=InnoDB AUTO_INCREMENT=306 DEFAULT CHARSET=utf8;

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
) ENGINE=InnoDB AUTO_INCREMENT=2211 DEFAULT CHARSET=utf8mb4;

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
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
