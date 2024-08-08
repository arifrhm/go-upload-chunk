import http from 'k6/http';
import { check, sleep } from 'k6';
import { Trend } from 'k6/metrics';

// Define the metrics
const uploadDuration = new Trend('upload_duration');

// File upload details
const filePath = './Head First Python  A Brain-Friendly Guide ( PDFDrive ).pdf'; // Path to your file
const uploadUrl = 'http://localhost:8001/upload-chunk'; // URL to your upload endpoint

// Chunk size in bytes (1MB)
const CHUNK_SIZE = 1024 * 1024; // 1MB

// Read the file into a buffer in the init stage
const file = open(filePath, 'b'); // 'b' for binary mode
const fileSize = file.byteLength;
console.log(fileSize);
const totalChunks = Math.ceil(fileSize / CHUNK_SIZE);

// Form data fields
const fileName = 'Head First Python A Brain-Friendly Guide ( PDFDrive ).pdf'; // Use the file name required by your server

export const options = {
  vus: 1, // Number of virtual users
  duration: '60s', // Test duration
};

export default function () {
  for (let i = 0; i < totalChunks; i++) {
    // Calculate the start and end of the chunk
    const start = i * CHUNK_SIZE;
    const end = Math.min(start + CHUNK_SIZE, fileSize);
    
    // Slice the file into chunks
    const chunk = file.slice(start, end);

    // Define the payload for the chunk upload
    const payload = {
      file: http.file(chunk, fileName),
      chunk_index: i.toString(),
      total_chunks: totalChunks.toString(),
      file_name: fileName,
    };

    // Measure the upload duration
    const startTime = new Date().getTime();
    const response = http.post(uploadUrl, payload);
    const endTime = new Date().getTime();

    // Log the duration
    const duration = endTime - startTime;
    uploadDuration.add(duration);

    // Check the response
    check(response, {
      'is status 200': (r) => r.status === 200,
    });
  }
}
