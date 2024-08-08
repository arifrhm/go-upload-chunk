const puppeteer = require("puppeteer");
const path = require("path");
const fs = require("fs");

(async () => {
  // Launch a headless browser
  const browser = await puppeteer.launch({ headless: true }); // Set headless to true for non-UI mode
  const page = await browser.newPage();

  // Open the HTML file or URL where the file upload is located
  await page.goto("http://localhost:8001/upload"); // Update with the path to your HTML file

  // Wait for the file input element to be present
  const fileInputSelector = "#fileInput";
  const fileInputElement = await page.waitForSelector(fileInputSelector);

  // Define the path to the file
  const filePath = path.resolve("./Head First Python  A Brain-Friendly Guide ( PDFDrive ).pdf"); // Update with the path to your file

  // Calculate the total file size and chunk size
  const fileSize = fs.statSync(filePath).size;
  const chunkSize = 1024 * 1024; // 1MB
  const totalChunks = Math.ceil(fileSize / chunkSize);

  // Use the uploadFile method to set the file on the input element
  await fileInputElement.uploadFile(filePath);

  // Variables to capture request IDs
  let requestIds = [];

  // Set up request interception to capture request IDs
  page.on('request', request => {
    if (request.url().includes("/upload-chunk")) {
      requestIds.push(request.id); // Capture request ID
      console.log(`Captured Request ID: ${request.id}`);
    }
  });

  // Click the upload button to start the upload process
  const uploadButtonSelector = "body > button:nth-child(2)";
  await page.waitForSelector(uploadButtonSelector); // Ensure the button is present
  await page.click(uploadButtonSelector);

  // Wait for the response for the last chunk upload
  await page.waitForResponse((response) => {
    const url = response.url();
    const status = response.status();
    const responseId = response.request().id;
    
    // Check if the response ID is the last captured ID
    const isLastRequest = requestIds.length == totalChunks && responseId === requestIds[requestIds.length - 1];
    return url.includes("/upload-chunk") &&
           isLastRequest &&
           status === 200;
  });

  // Close the browser
  await browser.close();
})();
