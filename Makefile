all:
	gb build

deploy: all
	git push -f heroku master
	heroku logs -t

createschema:
	heroku pg:psql < schema.sql

inspect:
	echo "select count(*) from logs; select * from logs order by id desc limit 15;" | heroku pg:psql
