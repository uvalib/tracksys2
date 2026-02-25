START TRANSACTION;

DROP TABLE iF EXISTS `notes_problems`;
DROP TABLE iF EXISTS `notes`;
DROP TABLE iF EXISTS `project_equipment`;
DROP TABLE iF EXISTS `assignments`;
DROP TABLE iF EXISTS `projects`;
DROP TABLE iF EXISTS `categories`;
DROP TABLE iF EXISTS `problems`;
DROP TABLE iF EXISTS `steps`;
DROP TABLE iF EXISTS `workstation_equipment`;
DROP TABLE iF EXISTS `workflows`;
DROP TABLE iF EXISTS `workstations`;
DROP TABLE iF EXISTS `equipment`;
DROP TABLE iF EXISTS `message_recipients`;
DROP TABLE iF EXISTS `messages`;

COMMIT;