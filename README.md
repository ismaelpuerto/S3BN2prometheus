# S3BN2prometheus

Generate endpoint prometheus (/metrics) when receive event http from Ceph Storage using bucket notification

Supported metrics (counter):
* s3_objectCreated_completeMultipartUpload
* s3_objectCreated_copy
* s3_objectCreated_post
* s3_objectCreated_put
* s3_objectCreated_wildcard
* s3_objectRemoved_delete
* s3_objectRemoved_deleteMarkerCreated
* s3_objectRemoved_wildcard
