UPDATE labels
SET
	value='mail'
WHERE
	value='message';

DELETE FROM labels
WHERE
	value='email';
