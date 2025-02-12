SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `couch_film_critic_db`
--

-- --------------------------------------------------------

--
-- Table structure for table `users`
--

CREATE TABLE IF NOT EXISTS `users` (
  `id` int  NOT NULL UNIQUE AUTO_INCREMENT,
  `nick_name` varchar(15) COLLATE utf8mb4_unicode_ci NOT NULL UNIQUE,
  `email` varchar(31) COLLATE utf8mb4_unicode_ci NOT NULL UNIQUE,
  `password_hash` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `signup_date` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
  
-- --------------------------------------------------------

--
-- Table structure for table `reviews`
--

CREATE TABLE IF NOT EXISTS `reviews` (
  `id` int NOT NULL UNIQUE AUTO_INCREMENT,
  `work_title` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `genres` varchar(31) COLLATE utf8mb4_unicode_ci NOT NULL,
  `work_type` varchar(12) COLLATE utf8mb4_unicode_ci NOT NULL,
  `review` varchar(400) COLLATE utf8mb4_unicode_ci NOT NULL,
  `rating` int COLLATE utf8mb4_unicode_ci NOT NULL,
  `create_date` datetime NOT NULL,
  `author` varchar(15) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

COMMIT;
/*
CREATE TABLE `reviews` (
  `id` int NOT NULL,
  `title` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `content` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `created` datetime NOT NULL,
  `expires` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
*/ 