```sh
# Register
curl -X POST https://<stagePlanFrom_TFoutput>.execute-api.<awsRegion>.amazonaws.com/blog/register -d '{"email": "foo@bar.test", "password": "changemeCHANGEME", "nombre": "foo", "apellidos": "bar"}'


# Login (get user token)
curl -X POST https://<stagePlanFrom_TFoutput>.execute-api.<awsRegion>.amazonaws.com/blog/login -d '{"email": "foo@bar.test", "password": "changemeCHANGEME", "nombre": "foo", "apellidos": "bar"}' | jq .

token=$(curl -X POST https://<stagePlanFrom_TFoutput>.execute-api.<awsRegion>.amazonaws.com/blog/login -d '{"email": "foo@bar.test", "password": "changemeCHANGEME", "nombre": "foo", "apellidos": "bar"}' | jq .token | sed s/\"//g)


# ViewProfile
curl -X GET https://<stagePlanFrom_TFoutput>.execute-api.<awsRegion>.amazonaws.com/blog/profile?id=65c779496a5ad8d4f812a128 -H 'Authorization: Bearer '$token


# ModifyProfile
curl -X PUT https://<stagePlanFrom_TFoutput>.execute-api.<awsRegion>.amazonaws.com/blog/modify?id=65c779496a5ad8d4f812a128 -H 'Authorization: Bearer '$token -d '{"email": "foo2@bar2.test", "password": "changemeCHANGEME", "nombre": "foo2", "apellidos": "bar2"}'


# SendTweet
curl -X POST https://<stagePlanFrom_TFoutput>.execute-api.<awsRegion>.amazonaws.com/blog/tweet -H 'Authorization: Bearer '$token -d '{"mensaje": "Este es el primer tweet"}'


# ReadTweets
# id == userId
curl -X GET "https://<stagePlanFrom_TFoutput>.execute-api.<awsRegion>.amazonaws.com/blog/tweet?id=65c779496a5ad8d4f812a128&page=1" -H 'Authorization: Bearer '$token


# DeleteTweet
# id == tweetId
curl -X DELETE https://<stagePlanFrom_TFoutput>.execute-api.<awsRegion>.amazonaws.com/blog/tweet?id=65c76a7c669f0d011929f6a4 -H 'Authorization: Bearer '$token


# Upload avatar
curl -X POST https://<stagePlanFrom_TFoutput>.execute-api.<awsRegion>.amazonaws.com/blog/avatar -H 'Authorization: Bearer '$token -F 'avatar=@/absolute/path/to/image.jpg'


# Upload banner
curl -X POST https://<stagePlanFrom_TFoutput>.execute-api.<awsRegion>.amazonaws.com/blog/banner -H 'Authorization: Bearer '$token -F 'banner=@/absolute/path/to/image.jpg'


# Get avatar
#curl -X GET https://<stagePlanFrom_TFoutput>.execute-api.<awsRegion>.amazonaws.com/blog/avatar?id=65c779496a5ad8d4f812a128 -H 'Authorization: Bearer '$token -v --output 000a.jpg


# Get banner
#curl -X GET https://<stagePlanFrom_TFoutput>.execute-api.<awsRegion>.amazonaws.com/blog/banner?id=65c779496a5ad8d4f812a128 -H 'Authorization: Bearer '$token -v --output 000b.jpg


# New relation
curl -X POST https://<stagePlanFrom_TFoutput>.execute-api.<awsRegion>.amazonaws.com/blog/relation?id=65c779496a5ad8d4f812a128 -H 'Authorization: Bearer '$token | jq .

# Get relation
curl -X GET https://<stagePlanFrom_TFoutput>.execute-api.<awsRegion>.amazonaws.com/blog/relation?id=65c779496a5ad8d4f812a128 -H 'Authorization: Bearer '$token | jq .

# Delete relation
curl -X DELETE https://<stagePlanFrom_TFoutput>.execute-api.<awsRegion>.amazonaws.com/blog/relation?id=65c779496a5ad8d4f812a128 -H 'Authorization: Bearer '$token | jq .


# Get users
# we need to create some users
curl -X POST https://<stagePlanFrom_TFoutput>.execute-api.<awsRegion>.amazonaws.com/blog/register -d '{"email": "foo2@bar2.test", "password": "changemeCHANGEME2", "nombre": "foo2", "apellidos": "bar2"}'
curl -X POST https://<stagePlanFrom_TFoutput>.execute-api.<awsRegion>.amazonaws.com/blog/register -d '{"email": "foo3@bar3.test", "password": "changemeCHANGEME3", "nombre": "foo3", "apellidos": "bar3"}'
curl -X POST https://<stagePlanFrom_TFoutput>.execute-api.<awsRegion>.amazonaws.com/blog/register -d '{"email": "foo4@bar4.test", "password": "changemeCHANGEME4", "nombre": "foo4", "apellidos": "bar4"}'
# we need to create some relations too
# user 1 follows users 2, 3, 4
curl -X POST https://<stagePlanFrom_TFoutput>.execute-api.<awsRegion>.amazonaws.com/blog/relation?id=65c779946a5ad8d4f812a12a -H 'Authorization: Bearer '$token
curl -X POST https://<stagePlanFrom_TFoutput>.execute-api.<awsRegion>.amazonaws.com/blog/relation?id=65c779b86a5ad8d4f812a12c -H 'Authorization: Bearer '$token
curl -X POST https://<stagePlanFrom_TFoutput>.execute-api.<awsRegion>.amazonaws.com/blog/relation?id=65c779c26a5ad8d4f812a12e -H 'Authorization: Bearer '$token
# user 2 follows user 1
token2=$(curl -X POST https://<stagePlanFrom_TFoutput>.execute-api.<awsRegion>.amazonaws.com/blog/login -d '{"email": "foo2@bar2.test", "password": "changemeCHANGEME2", "nombre": "foo2", "apellidos": "bar2"}' | jq .token | sed s/\"//g)
curl -X POST https://<stagePlanFrom_TFoutput>.execute-api.<awsRegion>.amazonaws.com/blog/relation?id=65c779496a5ad8d4f812a128 -H 'Authorization: Bearer '$token2
# user 1 following
curl -X GET 'https://<stagePlanFrom_TFoutput>.execute-api.<awsRegion>.amazonaws.com/blog/users?page=1&type=follow&search' -H 'Authorization: Bearer '$token | jq .
# user 2 following
curl -X GET 'https://<stagePlanFrom_TFoutput>.execute-api.<awsRegion>.amazonaws.com/blog/users?page=1&type=follow&search' -H 'Authorization: Bearer '$token2 | jq .
# null -> user 1 already follow user 2, 3, 4
curl -X GET 'https://<stagePlanFrom_TFoutput>.execute-api.<awsRegion>.amazonaws.com/blog/users?page=1&type=new&search' -H 'Authorization: Bearer '$token | jq .
# user 2 already follow user 1 (it does not apear) but he can follow user 3 and 4
curl -X GET 'https://<stagePlanFrom_TFoutput>.execute-api.<awsRegion>.amazonaws.com/blog/users?page=1&type=new&search' -H 'Authorization: Bearer '$token2 | jq .


# followersTweets
# user 2 follows user 1, lets view the tweets published by this user
curl -X GET 'https://<stagePlanFrom_TFoutput>.execute-api.<awsRegion>.amazonaws.com/blog/followersTweets?page=1' -H 'Authorization: Bearer '$token2 | jq .
```
