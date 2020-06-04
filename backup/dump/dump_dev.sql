-- MySQL dump 10.13  Distrib 5.7.30, for Linux (x86_64)
--
-- Host: 127.0.0.1    Database: webservice
-- ------------------------------------------------------
-- Server version	5.7.29

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `adeslas`
--

DROP TABLE IF EXISTS `adeslas`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `adeslas` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `lea_id` int(10) unsigned DEFAULT NULL,
  `product` varchar(255) DEFAULT NULL,
  `landing` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_adeslas_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `alterna`
--

DROP TABLE IF EXISTS `alterna`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `alterna` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `lea_id` int(10) unsigned DEFAULT NULL,
  `install_type` varchar(255) DEFAULT NULL,
  `cpu_s` varchar(255) DEFAULT NULL,
  `street` varchar(255) DEFAULT NULL,
  `number` varchar(255) DEFAULT NULL,
  `postal_code` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_alterna_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `creditea`
--

DROP TABLE IF EXISTS `creditea`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `creditea` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `lea_id` int(10) unsigned DEFAULT NULL,
  `requested_amount` varchar(255) DEFAULT NULL,
  `contract_type` varchar(255) DEFAULT NULL,
  `net_income` varchar(255) DEFAULT NULL,
  `out_of_schedule` varchar(255) DEFAULT NULL,
  `asnef` tinyint(1) DEFAULT NULL,
  `already_client` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_creditea_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `kinkon`
--

DROP TABLE IF EXISTS `kinkon`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `kinkon` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `lea_id` int(10) unsigned DEFAULT NULL,
  `coverture` varchar(255) DEFAULT NULL,
  `state` varchar(255) DEFAULT NULL,
  `town` varchar(255) DEFAULT NULL,
  `street` varchar(255) DEFAULT NULL,
  `number` varchar(255) DEFAULT NULL,
  `floor` varchar(255) DEFAULT NULL,
  `cov_phone` varchar(255) DEFAULT NULL,
  `phone` varchar(255) DEFAULT NULL,
  `phone_provider` varchar(255) DEFAULT NULL,
  `mobile_phone` varchar(255) DEFAULT NULL,
  `mobile_phone_provider` varchar(255) DEFAULT NULL,
  `mobile_phone2` varchar(255) DEFAULT NULL,
  `mobile_phone_provider2` varchar(255) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `surname` varchar(255) DEFAULT NULL,
  `idnumber` varchar(255) DEFAULT NULL,
  `mail` varchar(255) DEFAULT NULL,
  `contact_phone` varchar(255) DEFAULT NULL,
  `account_holder` varchar(255) DEFAULT NULL,
  `account_number` varchar(255) DEFAULT NULL,
  `product` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_kinkon_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `leads`
--

DROP TABLE IF EXISTS `leads`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `leads` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `legacy_id` bigint(20) DEFAULT NULL,
  `lea_ts` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `lea_smartcenter_id` varchar(255) DEFAULT NULL,
  `passport_id` varchar(255) DEFAULT NULL,
  `passport_id_grp` varchar(255) DEFAULT NULL,
  `sou_id` bigint(20) DEFAULT NULL,
  `leatype_id` bigint(20) DEFAULT NULL,
  `utm_source` varchar(255) DEFAULT NULL,
  `sub_source` varchar(255) DEFAULT NULL,
  `lea_phone` varchar(255) DEFAULT NULL,
  `lea_mail` varchar(255) DEFAULT NULL,
  `lea_name` varchar(255) DEFAULT NULL,
  `lea_dni` varchar(255) DEFAULT NULL,
  `lea_url` text,
  `lea_ip` varchar(255) DEFAULT NULL,
  `ga_client_id` varchar(255) DEFAULT NULL,
  `is_smart_center` tinyint(1) DEFAULT NULL,
  `gclid` varchar(255) DEFAULT NULL,
  `domain` varchar(255) DEFAULT NULL,
  `observations` text,
  PRIMARY KEY (`id`),
  KEY `idx_leads_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=48 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `leadtypes`
--

DROP TABLE IF EXISTS `leadtypes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `leadtypes` (
  `leatype_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `leatype_description` varchar(255) DEFAULT NULL,
  `leatype_idcrm` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`leatype_id`)
) ENGINE=InnoDB AUTO_INCREMENT=31 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `microsoft`
--

DROP TABLE IF EXISTS `microsoft`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `microsoft` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `lea_id` int(10) unsigned DEFAULT NULL,
  `computer_type` varchar(255) DEFAULT NULL,
  `sector` varchar(255) DEFAULT NULL,
  `budget` varchar(255) DEFAULT NULL,
  `performance` varchar(255) DEFAULT NULL,
  `movility` varchar(255) DEFAULT NULL,
  `office365` varchar(255) DEFAULT NULL,
  `usecase` varchar(255) DEFAULT NULL,
  `comments` varchar(255) DEFAULT NULL,
  `product_type` varchar(255) DEFAULT NULL,
  `product_name` varchar(255) DEFAULT NULL,
  `product_id` varchar(255) DEFAULT NULL,
  `original_price` varchar(255) DEFAULT NULL,
  `price` varchar(255) DEFAULT NULL,
  `brand` varchar(255) DEFAULT NULL,
  `discount_percentage` varchar(255) DEFAULT NULL,
  `discount_code` varchar(255) DEFAULT NULL,
  `processor_type` varchar(255) DEFAULT NULL,
  `disk_capacity` varchar(255) DEFAULT NULL,
  `graphics` varchar(255) DEFAULT NULL,
  `wireless_interface` varchar(255) DEFAULT NULL,
  `devices_average_age` varchar(255) DEFAULT NULL,
  `devices_operating_system` varchar(255) DEFAULT NULL,
  `devices_hang_frequency` varchar(255) DEFAULT NULL,
  `devices_number` varchar(255) DEFAULT NULL,
  `devices_last_year_repairs` varchar(255) DEFAULT NULL,
  `devices_startup_time` varchar(255) DEFAULT NULL,
  `pageindex` tinyint(1) DEFAULT NULL,
  `oldsouid` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_microsoft_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `person`
--

DROP TABLE IF EXISTS `person`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `person` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `project`
--

DROP TABLE IF EXISTS `project`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `project` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `rcableexp`
--

DROP TABLE IF EXISTS `rcableexp`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rcableexp` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `lea_id` int(10) unsigned DEFAULT NULL,
  `location` varchar(255) DEFAULT NULL,
  `answer` varchar(255) DEFAULT NULL,
  `resp_values` varchar(255) DEFAULT NULL,
  `coverture` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_rcableexp_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `sources`
--

DROP TABLE IF EXISTS `sources`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `sources` (
  `sou_id` bigint(20) DEFAULT NULL,
  `sou_description` varchar(255) DEFAULT NULL,
  `sou_idcrm` bigint(20) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `status`
--

DROP TABLE IF EXISTS `status`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `status` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `untraceable`
--

DROP TABLE IF EXISTS `untraceable`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `untraceable` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `lea_id` bigint(20) DEFAULT NULL,
  `phone` varchar(255) DEFAULT NULL,
  `sou_id` bigint(20) DEFAULT NULL,
  `sms_id` int(10) unsigned DEFAULT NULL,
  `sms_date` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `untraceable_test`
--

DROP TABLE IF EXISTS `untraceable_test`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `untraceable_test` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `lea_id` bigint(20) DEFAULT NULL,
  `phone` varchar(255) DEFAULT NULL,
  `sou_id` bigint(20) DEFAULT NULL,
  `ddi` varchar(255) DEFAULT NULL,
  `sms_date` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `work`
--

DROP TABLE IF EXISTS `work`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `work` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `person_id` int(10) unsigned DEFAULT NULL,
  `project_id` int(10) unsigned DEFAULT NULL,
  `status_id` int(10) unsigned DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `estimated_time` double DEFAULT NULL,
  `real_time` double DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_work_deleted_at` (`deleted_at`),
  KEY `work_person_id_person_id_foreign` (`person_id`),
  KEY `work_project_id_project_id_foreign` (`project_id`),
  KEY `work_status_id_status_id_foreign` (`status_id`),
  CONSTRAINT `work_person_id_person_id_foreign` FOREIGN KEY (`person_id`) REFERENCES `person` (`id`),
  CONSTRAINT `work_project_id_project_id_foreign` FOREIGN KEY (`project_id`) REFERENCES `project` (`id`),
  CONSTRAINT `work_status_id_status_id_foreign` FOREIGN KEY (`status_id`) REFERENCES `status` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2020-06-03 17:39:01
