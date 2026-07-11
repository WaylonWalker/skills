---
name: dropper
description: Upload files to Dropper with curl. Use this skill when the user wants to send a local file to https://dropper.waylonwalker.com/upload, check the response, or debug a failed upload. This should trigger whenever the task involves Dropper uploads, even if the user only mentions a file path and not the service name.
---

# Dropper Upload

Use this skill to upload a local file to Dropper and verify the response.

## Basic upload

Use `curl` with multipart form upload:

```bash
curl -X POST https://dropper.waylonwalker.com/upload \
  -F "file=@/path/to/image.png"
```

## Response handling

Treat the response as part of the task.

- Print the response so the user can see the uploaded file URL or server message.
- If the upload fails, surface the HTTP status and response body.
- If the user gives a different file path, substitute it directly.

## If the upload fails

Check the obvious causes first:

- The file path is wrong.
- The file does not exist.
- The server returns a non-2xx status.
- The file is too large or the server rejects the content type.

Keep the fix minimal. Use the exact upload endpoint unless the user gives a different one.
