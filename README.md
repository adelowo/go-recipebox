### Implementation of Freecodecamp RecipeBoard in Golang

I currently am learning Go and i decided to make use of this project from freecodecamp as it was the next thing on my todo - albeit in Javascript.


User story

- Users can signup.
- Users can login.
- Users can add recipes.
- Users can view recipes
- Users can delete recipes. This is limited to recipes they created. They cannot delete recipes created by another user.
- Users can edit recipes.


### Setup

```bash
$ cat schema.sql | sqlite3 recipebox.db
#if you decide to use another db source other than recipebox.db, update the database key in the config/config.json file.
```


> There currently are no integration tests for this but there are unit tests for wrappers i added. Haven't wrapped my head around the httptest package though.
