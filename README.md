# A+ Project Web Site - [Link](https://www.youtube.com/watch?reload=9&reload=9&v=sudUExRVw1I)

## Description

- A+ Project is a educational website used by professors and students.

- There are three sections of the website.

1) Admins section have almost all privilages like (adding professors, students on the system, inserting new courses, update them and insert other admins on it )

2) Professors section that help professors to set exams, edit exams, enable exam for students in any time, see degrees for students and can edit in any degree for any student.

3) Students section where student login, enter exam and see degree for his/her exam

## Structure

Each table in project consists of three files

1) Model File
 represent the the table field ,keys and type of each one

2) Url File
 represent The End Points of the table and The method type of Each one like (POST,GET,PUT,DELETE) and functions on API File

3) API File
represent The Caller Function of each End Point from Url File


- This is my project structure

    ```
    └── $HOME
        └── go
            ├── bin
            ├── pkg
            └── src
                └── gitlab.com
                    └── mohamedzaky
                        └── aplusProject
    ```

## Implmentaion

- Back End Language  [GoLang](https://golang.org/)
- Framework and Middleware [Echo](https://echo.labstack.com/)

- ORM library [Gorm](https://github.com/jinzhu/gorm)

- Database [PostgreSQL](https://www.postgresql.org/)

## Installation
Clone the project:
```
    git clone git@github.com:mohammedzaky/AplusProject.git
```

## Run Database
- In folder configration there is a db.yaml where is a configration of Database (set the configration of your database)

- Must make this configration to run the DataBase or Edit the db.yaml file for custom purpose  
```
    Host: localhost
    Port: 5432
    User: postgres 
    Dbname: postgres
    Password: test1234 
    Sslmode: disable
    ResetDB: false
```

- ResetDB if there is

    1) True then all tables will be truncated & dropped and recreated again

    2) False then will continue to use database without dropped any table
