CREATE DATABASE sql_guestbook; 
GO
USE sql_guestbook;
GO
CREATE TABLE guestlog (entrydate DATETIME, name NVARCHAR(30), phone NVARCHAR(30), message TEXT, sentiment_score NVARCHAR(30)); 
GO
INSERT INTO guestlog VALUES ('2017-5-2 23:59:59', 'anonymous', '12158379120', 'Get busy living, or get busy dying', '0.9950121'); 
GO
INSERT INTO guestlog VALUES ('2017-4-15 23:59:59', 'anonymous', '19192310925', 'That rug really tied the room together', '0.6625549');
GO