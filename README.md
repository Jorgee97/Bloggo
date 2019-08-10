# Bloggo
A little Blog API written in Go

# Features
- Register/Login system
- Authentication using JWT
- Storage with MongoDB

# Endpoints

| HTTP Method | URI path                 | Description                                                    |
|-------------|--------------------------|----------------------------------------------------------------|
| POST        | /signup                  | Allow the client to create a new account                       |
| POST        | /login                   | Allow the client to login with credentials                     |
| GET         | /blog/username/:username | Retrieve the posts that are not private from a respective user |
| GET         | /blog/:id                | Retrieve a post by ID                                          |
| PUT         | /blog/:id                | Update a respective post by its ID                             |
| DELETE      | /blog/:id                | Delete a blog by its respective ID                             |
| POST        | /blog                    | Create a new post                                              |
| GET         | /blog                    | Retrieves all posts from all users that are not private        |

# Notes on the implementation
1. This is not a production ready blog service, it may contain bugs and as it is my first go project it may not be well structured.
2. To be able to run this application you need to install and configure your MongoDB server

# License
```
MIT License

Copyright (c) [year] [fullname]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
