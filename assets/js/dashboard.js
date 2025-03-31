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

  /** @type HTMLButtonElement | null */
  weights_button = null;

  /** @type HTMLButtonElement | null */
  dnr_dpmo_button = null;

  constructor() {
    this.pdf_drag_zone = document.getElementById("pdf_drag_zone");
    this.pdf_input = document.getElementById("dropzone-file");
    this.upload_list_body = document.getElementById("upload_list_body");
    this.weights_button = document.getElementById("weights_button");

    this.form_container = document.getElementById("form_container");

    this.dnr_dpmo_button = document.getElementById("station_dropdown");
    this.station_dropdown = document.getElementById("dropdown");

    this.stations_form = document.getElementById("station_form");
    this.stations_form.addEventListener(
      "submit",
      this.submit_station.bind(this),
    );

    this.current_station = null;

    this.is_displaying_station_form = false;
    this.is_displaying_weights_form = false;

    if (
      !this.pdf_drag_zone ||
      !this.pdf_input ||
      !this.upload_list_body ||
      !this.weights_button ||
      !this.dnr_dpmo_button
    ) {
      throw new Error("There is no drag zone");
    }

    this.pdf_drag_zone.addEventListener("drop", this.handle_drop.bind(this));
    this.pdf_drag_zone.addEventListener(
      "dragover",
      this.handle_drag_over.bind(this),
    );
    this.dnr_dpmo_button.addEventListener(
      "click",
      this.toggle_dnr_dmpl.bind(this),
    );
    this.weights_button.addEventListener(
      "click",
      this.toggle_weights_botton.bind(this),
    );
  }

  /**
   *
   * Make the station form visible
   *
   */
  async toggle_dnr_dmpl(e) {
    this.is_displaying_station_form = !this.is_displaying_station_form;

    if (this.is_displaying_station_form) {
      // Just let the error break things, f - try/catch
      if (!e.target.id) {
        throw new Error("there is no id!");
      }

      const station_data = await this.get_station_data();
      const j_station_data = await station_data.json();

      let station_to_display = j_station_data.filter(
        (el) => el.station === e.target.id,
      );

      if (station_to_display.length < 1) {
        throw new Error("there is no station to display");
      }

      // Make the form show the current values
      station_to_display = station_to_display[0];
      const inputs = this.stations_form.querySelectorAll("input");
      for (const el of inputs) {
        if (station_to_display[el.id]) {
          el.value = station_to_display[el.id];
        }
      }

      this.current_station = station_to_display;

      this.form_container.style.display = "flex";
      this.stations_form.style.display = "block";
      this.station_dropdown.classList.add("hidden");

      return;
    }

    this.stations_form.style.display = "none";
    this.form_container.style.display = "none";
  }

  async submit_station(e) {
    e.preventDefault();
    const target = e.target;

    if (!target) {
      throw new Error("there is no form target");
    }

    const inputs = target.querySelectorAll("input");

    for (const el of inputs) {
      this.current_station[el.name] = parseInt(el.value);
    }

    const res = await this.update_station();
    console.log("what ia the res: ", res);

    this.current_station = null;

    this.toggle_dnr_dmpl();
  }

  async update_station() {
    const endpoint = "/station";
    let current_domain = window.location.href;

    current_domain = current_domain.replace("dashboard", "api");

    const json_station = JSON.stringify(this.current_station);

    return fetch(current_domain + endpoint, {
      method: "POST",
      body: json_station,
      headers: {
        "Content-Type": "application/json",
      },
    });
  }

  toggle_weights_botton(e) {
    console.log("the toggle dnr was clicked");
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

        const new_name = name.replace("\.pdf", "\.csv");

        const blob = new Blob(body.csv, { type: "text/plain" });
        const download_url = URL.createObjectURL(blob);

        const link = document.createElement("a");
        link.href = download_url;
        link.download = new_name;
        link.textContent = "download";

        download.append(link);
        // link.click();

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

  async get_station_data() {
    const endpoint = "/station";
    let current_domain = window.location.href;

    current_domain = current_domain.replace("dashboard", "api");

    return fetch(current_domain + endpoint, {
      method: "GET",
    });
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
    const endpoint = "/pdf_upload";
    let current_domain = window.location.href;

    current_domain = current_domain.replace("dashboard", "api");

    const form_data = new FormData();
    form_data.append("file", file);

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
