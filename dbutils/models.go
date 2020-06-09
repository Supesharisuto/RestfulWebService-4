package dbutils

const video = `
	CREATE TABLE IF NOT EXISTS video (
           ID INTEGER PRIMARY KEY AUTOINCREMENT,
           NAME VARCHAR(64) NULL,
           ASSETID VARCHAR(64) NULL
        )
`

const assets = `
	CREATE TABLE IF NOT EXISTS assets (
          ID INTEGER PRIMARY KEY AUTOINCREMENT,
          NAME VARCHAR(64) NULL,
          CREATE_TIME TIME NULL
        )
`
const schedule = `
	CREATE TABLE IF NOT EXISTS schedule (
	  ID INTEGER PRIMARY KEY AUTOINCREMENT,
          VIDEO_ID INT,
          ASSETS_ID INT,
          SCHED_DELETE_TIME TIME,
          FOREIGN KEY (VIDEO_ID) REFERENCES video(ID),
          FOREIGN KEY (ASSETS_ID) REFERENCES assets(ID)
        )
`
