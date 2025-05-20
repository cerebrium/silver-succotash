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

    this.tier_button = document.getElementById("tier_dropdown");
    this.tier_dropdown = document.getElementById("dropdown_tier");

    this.stations_form = document.getElementById("station_form");
    this.stations_form.addEventListener(
      "submit",
      this.submit_station.bind(this),
    );

    this.weights_form = document.getElementById("weights_form");
    this.weights_form.addEventListener(
      "submit",
      this.submit_weights.bind(this),
    );

    this.tiers_form = document.getElementById("tier_form");
    this.tiers_form.addEventListener("submit", this.submit_tiers.bind(this));

    this.current_station = null;
    this.current_tier = null;

    this.is_displaying_station_form = false;
    this.is_displaying_weights_form = false;
    this.is_displaying_tier_form = false;

    if (
      !this.pdf_drag_zone ||
      !this.pdf_input ||
      !this.upload_list_body ||
      !this.weights_button ||
      !this.dnr_dpmo_button ||
      !this.tier_button
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

    this.tier_button.addEventListener(
      "click",
      this.toggle_tier_form.bind(this),
    );
    this.weights_button.addEventListener(
      "click",
      this.toggle_weights_botton.bind(this),
    );

    this.getUserEmail();
  }

  // Example usage (replace with your actual JWT):
  async getUserEmail() {
    const user_display = document.querySelector(".cl-userButtonTrigger");

    if (!user_display) {
      return setTimeout(() => {
        this.getUserEmail();
      }, 10);
    }

    user_display.click();

    this.findEmail();
  }

  async findEmail() {
    const email_box = document.querySelector(
      ".cl-userPreviewSecondaryIdentifier",
    );

    if (!email_box) {
      return setTimeout(() => {
        this.findEmail();
      }, 10);
    }

    if (
      (email_box.textContent &&
        email_box.textContent === "nicholas.m.shankland@gmail.com") ||
      email_box.textContent === "r.marconi@h2ologistics.co.uk"
    ) {
      this.is_authed = true;
    } else {
      const dropdown_button = document.getElementById("dropdownDefaultButton");

      if (!dropdown_button) {
        throw new Error("this is busted, no dropdown");
      }
      dropdown_button.style.display = "none";
      this.weights_button.style.display = "none";
      this.is_authed = false;
    }

    const user_display = document.querySelector(".cl-userButtonTrigger");

    user_display.click();
  }

  // TODO: Implement update
  async submit_tiers(e) {
    e.preventDefault();
    const target = e.target;

    if (!target) {
      throw new Error("there is no form target");
    }

    const inputs = target.querySelectorAll("input");
    const tiers = {};

    for (const el of inputs) {
      switch (el.name) {
        case "fantastic_plus":
          tiers["fan_plus"] = parseFloat(parseFloat(el.value).toFixed(2));
          continue;
        case "fantastic":
          tiers["fan"] = parseFloat(parseFloat(el.value).toFixed(2));
          continue;

        default:
          tiers[el.name] = parseFloat(parseFloat(el.value).toFixed(2));
      }
    }

    // Get all the tiers

    const tier_data = await this.get_tier_data();
    const j_tier_data = await tier_data.json();

    // find the one to update, then push the changes, send all the
    // data as json

    for (let vals of j_tier_data) {
      if (vals.name.toLowerCase() === this.current_tier.toLowerCase()) {
        for (const key in vals) {
          if (key === "id" || key === "name") {
            continue;
          }
          vals[key] = tiers[key];
        }
      }
    }

    const endpoint = "/tiers";
    let current_domain = window.location.href;

    current_domain = current_domain.replace("dashboard", "api");

    const json_tier_data = JSON.stringify(j_tier_data);
    this.current_tier = null;

    fetch(current_domain + endpoint, {
      method: "POST",
      body: json_tier_data,
      headers: {
        "Content-Type": "application/json",
      },
    });

    this.tiers_form.style.display = "none";
    this.form_container.style.display = "none";

    return;
  }

  async submit_weights(e) {
    e.preventDefault();
    const target = e.target;

    if (!target) {
      throw new Error("there is no form target");
    }

    const inputs = target.querySelectorAll("input");

    const weights = {};

    for (const el of inputs) {
      weights[el.name] = parseInt(el.value);
    }

    let current_sum = 0;

    for (const val of Object.values(weights)) {
      current_sum += val;
    }

    if (100 - current_sum > 3 || 100 - current_sum < -3) {
      Toastify({
        text: "Values are not close enough to 100 summed",
        duration: 3000,
        destination: "https://github.com/apvarun/toastify-js",
        newWindow: true,
        close: true,
        gravity: "top",
        position: "right",
        stopOnFocus: true,
        style: {
          background: "linear-gradient(to right, #FA8072, #CD5C5C)",
        },
        onClick: function () {}, // Callback after click
      }).showToast();

      // TODO: make a toast
    } else {
      await this.update_weights(weights);

      this.toggle_weights_botton();
    }
  }

  async update_weights(weights) {
    const endpoint = "/weights";
    let current_domain = window.location.href;

    current_domain = current_domain.replace("dashboard", "api");

    for (const [key, value] of Object.entries(weights)) {
      const new_val = parseInt(value) / 100;
      weights[key] = parseFloat(new_val.toFixed(3));
    }

    weights.ID = 1;

    const json_station = JSON.stringify(weights);

    return fetch(current_domain + endpoint, {
      method: "POST",
      body: json_station,
      headers: {
        "Content-Type": "application/json",
      },
    });
  }

  async toggle_weights_botton(e) {
    this.is_displaying_weights_form = !this.is_displaying_weights_form;

    if (this.is_displaying_weights_form) {
      const weights = await this.get_weights();
      const j_weights = await weights.json();

      const form_inputs = this.weights_form.querySelectorAll("input");

      for (const el of form_inputs) {
        if (j_weights[el.id]) {
          el.value = parseFloat(j_weights[el.id] * 100).toFixed(1);
        }
      }

      this.form_container.style.display = "flex";
      this.weights_form.style.display = "block";
      return;
    }

    this.weights_form.style.display = "none";
    this.form_container.style.display = "none";
  }

  async get_weights() {
    const endpoint = "/weights";
    let current_domain = window.location.href;

    current_domain = current_domain.replace("dashboard", "api");

    return fetch(current_domain + endpoint, {
      method: "GET",
    });
  }

  async get_tier_data() {
    const endpoint = "/tiers";
    let current_domain = window.location.href;

    current_domain = current_domain.replace("dashboard", "api");

    return fetch(current_domain + endpoint, {
      method: "GET",
    });
  }

  /**
   *
   * Make the tier form visible
   *
   */

  async toggle_tier_form(e) {
    if (!e.target.id) {
      throw new Error("there is no id!");
    }

    this.current_tier = e.target.id;

    const tier_data = await this.get_tier_data();
    const j_tier_data = await tier_data.json();

    let data_to_populate = null;
    for (const vals of j_tier_data) {
      if (vals.name.toLowerCase() === e.target.id.toLowerCase()) {
        data_to_populate = vals;
      }
    }

    const inputs = this.tiers_form.querySelectorAll("input");
    for (const el of inputs) {
      switch (el.id) {
        case "fantastic": {
          el.value = data_to_populate["fan"];
          continue;
        }
        case "fantastic_plus": {
          el.value = data_to_populate["fan_plus"];
          continue;
        }
        default:
          el.value = data_to_populate[el.id];
      }
    }

    this.form_container.style.display = "flex";
    this.tiers_form.style.display = "block";
    this.tier_dropdown.classList.add("hidden");

    return;
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

    await this.update_station();

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
    current_domain = current_domain.replace("/v1", "");

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
