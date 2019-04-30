CREATE TABLE IF NOT EXISTS `skel`.`user` (
  `user_id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(45) NULL,
  `email` VARCHAR(50) NULL,
  `password` VARCHAR(128) NULL,
  `active` TINYINT(1) NULL DEFAULT 1,
  PRIMARY KEY (`user_id`),
  UNIQUE INDEX `email_UNIQUE` (`email` ASC))
ENGINE = InnoDB