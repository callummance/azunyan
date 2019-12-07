db.createUser(
	{
	user: "admin",
	pwd: "admin",
	roles : [
		{
			role: "readWrite",
			db: "azunyan",
		}
	]}
)	
db.createCollection("songs");
db.songs.createIndex( {"title": "text", "artist": "text"}, {"default_language": "en", "language_override": "en"});
