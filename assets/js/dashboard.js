"use strict";

class UploadFile {
  /** @type Array<string> */
  uploaded_elements = [];

  /** @type HTMLElement | null */
  pdf_drag_zone = null;

  /** @type HTMLInputElement | null */
  pdf_input = null;

  /** @type boolean */
  is_dragging = false;

  constructor() {
    this.pdf_drag_zone = document.getElementById("pdf_drag_zone");
    this.pdf_input = document.getElementById("dropzone-file");
    this.upload_list_body = document.getElementById("upload_list_body");

    if (!this.pdf_drag_zone || !this.pdf_input || !this.upload_list_body) {
      throw new Error("There is no drag zone");
    }

    this.pdf_drag_zone.addEventListener("drop", this.handle_drop.bind(this));
    this.pdf_drag_zone.addEventListener(
      "dragover",
      this.handle_drag_over.bind(this),
    );
  }

  handle_drop(e) {
    e.preventDefault();
    const files = e.dataTransfer.files;

    if (!files.length) {
      alert("There was no file dropped!");
      return;
    }

    const name = files[0].name;
    if (this.uploaded_elements.includes(name)) {
      alert("File already dropped");
      return;
    }

    /*
     *
     * We want to:
     *
     * Update the file list to have a new entry
     * send a request to the backend with the pdf
     * create a spinner for the upload
     *
     */
    const [status, download] = this.handle_update_ui(name);

    this.handle_formatting_file(files[0]).then(async (res) => {
      /*
       *
       * This should spit out a csv file, which will need to go
       * into the just created elements download with id = name
       *
       */
      if (res.body) {
        const body = await res.json();

        /*
         *
         * Steps to download an object:
         * 1. create a hyperlink
         * 2. blobify the data
         * 3. set the content into the hyperlink
         * 4. click the hyperlink
         * 5. destroy the hyperlink
         *
         */

        const blob = new Blob(body.csv, { type: "text/plain" });
        const download_url = URL.createObjectURL(blob);

        const link = document.createElement("a");
        link.href = download_url;
        link.download = name;
        link.textContent = "download";

        download.append(link);
        link.click();

        status.textContent = "success";
      }
    });
  }

  /**
   * Updates the UI with a new row in the upload list table.
   *
   * @param {string} name - The name of the file to display in the UI.
   * @returns {void}
   */
  handle_update_ui(name) {
    const new_file_entry = document.createElement("tr");
    new_file_entry.classList.add(
      ...[
        "bg-white",
        "border-b",
        "dark:bg-gray-800",
        "dark:border-gray-700",
        "hover:bg-gray-50",
        "dark:hover:bg-gray-600",
        "h-14",
      ],
    );

    const title = document.createElement("th");
    title.scope = "row";
    title.textContent = name;
    title.classList.add(
      ...[
        "px-6",
        "py-4",
        "font-medium",
        "text-gray-900",
        "whitespace-nowrap",
        "dark:text-white",
      ],
    );

    const uploaded = document.createElement("td");
    uploaded.textContent = "success";
    uploaded.classList.add(...["px-6", "py-4"]);

    const status = document.createElement("td");
    status.textContent = "pending";
    status.classList.add(...["px-6", "py-4"]);

    const download = document.createElement("td");
    download.id = name;
    download.classList.add("download");

    new_file_entry.append(title);
    new_file_entry.append(uploaded);
    new_file_entry.append(status);
    new_file_entry.append(download);

    this.upload_list_body.append(new_file_entry);

    return [status, download];
  }

  /**
   * Processes a PDF file upload and performs specific actions (e.g., validation, rendering, or uploading).
   *
   * @param {e}
   * @param {file} file - the files metadata.
   * @returns {void}
   *
   * @throws {Error}
   */
  async handle_formatting_file(file) {
    const reader = new FileReader();

    reader.onload = async (e) => {
      const base64PDF = e.target.result.split(",")[1];
    };

    reader.readAsDataURL(file);
    let current_domain = window.location.href;
    current_domain = current_domain.replace("dashboard", "api");

    const endpoint = "/pdf_upload";

    const form_data = new FormData();

    form_data.file = file;

    console.log("form_data: ", form_data);
    return fetch(current_domain + endpoint, {
      method: "POST",
      body: form_data,
    });
  }

  async send_pdf() {}

  handle_drag_over(e) {
    e.preventDefault();
  }
}

new UploadFile();
