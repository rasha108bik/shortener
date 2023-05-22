/*
Package handlers provides all public API methods

CreateShortLink create short link and save in DB
CreateShorten create short link and save in DB with ReqCreateShorten model
ErrorHandler method which not include API methods
FetchURLs get all URLs from DB
GetOriginalURL get URL which include in url param
Ping is handler which check liveness server
ShortenBatch save URLs which include ReqShortenBatch model

Handlers all API public methods
interface.

NewHandler returns a newly initialized handler objects that implements the Handlers interface
interface.

DeleteURLs delete urls which include in reqBody
*/
package handlers
