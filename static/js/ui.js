/**
 * Small behaviors for the blog: local dates, the photo lightbox,
 * dropdown menus, the compose dialog's photo picker, and dropping a
 * freshly created post into the feed.
 */

const MAX_PHOTOS = 5;

export function initUI() {
    localizeDates(document);

    document.addEventListener('htmx:afterSwap', (event) => {
        localizeDates(event.detail.target);
        if (event.detail.target.id === 'form-response') onPosted(event.detail.target);
    });

    document.addEventListener('click', (event) => {
        // Close any open menu when clicking outside it.
        document.querySelectorAll('details.menu[open]').forEach((menu) => {
            if (!menu.contains(event.target)) menu.open = false;
        });

        const photo = event.target.closest('.photo');
        if (photo) openLightbox(photo);
    });

    document.addEventListener('keydown', (event) => {
        if (event.key !== 'Escape') return;
        document.querySelectorAll('details.menu[open]').forEach((menu) => (menu.open = false));
    });

    const lightbox = document.getElementById('image_modal');
    if (lightbox) {
        // Any click on the lightbox (image or backdrop) closes it.
        lightbox.addEventListener('click', () => lightbox.close());
        lightbox.addEventListener('close', () => (document.body.style.overflow = ''));
    }

    document.querySelectorAll('[data-dropzone]').forEach(initDropzone);
}

/** Rewrite server-rendered UTC dates as the reader's local date; the full time sits in the tooltip. */
function localizeDates(root) {
    root.querySelectorAll('time[data-local]').forEach((el) => {
        const date = new Date(el.getAttribute('datetime'));
        if (isNaN(date.getTime())) return;
        const pad = (n) => String(n).padStart(2, '0');
        el.textContent = `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}`;
        el.title = date.toLocaleString();
    });
}

function openLightbox(photo) {
    const lightbox = document.getElementById('image_modal');
    const image = document.getElementById('modal_image');
    if (!lightbox || !image) return;
    image.src = photo.currentSrc || photo.src;
    image.alt = photo.alt;
    lightbox.showModal();
    document.body.style.overflow = 'hidden';
}

function initDropzone(zone) {
    const input = zone.querySelector('input[type=file]');
    const preview = document.getElementById(input.dataset.preview);

    input.addEventListener('change', () => renderPreview(input, preview));

    zone.addEventListener('dragover', (event) => {
        event.preventDefault();
        zone.classList.add('is-over');
    });
    zone.addEventListener('dragleave', () => zone.classList.remove('is-over'));
    zone.addEventListener('drop', (event) => {
        event.preventDefault();
        zone.classList.remove('is-over');
        const files = [...event.dataTransfer.files].filter((file) => file.type.startsWith('image/'));
        if (!files.length) return;
        const transfer = new DataTransfer();
        files.slice(0, MAX_PHOTOS).forEach((file) => transfer.items.add(file));
        input.files = transfer.files;
        input.dispatchEvent(new Event('change'));
    });
}

function renderPreview(input, preview) {
    preview.innerHTML = '';
    const files = [...input.files].slice(0, MAX_PHOTOS);
    preview.classList.toggle('hidden', files.length === 0);

    files.forEach((file) => {
        const img = document.createElement('img');
        img.className = 'aspect-square w-full rounded-[3px] bg-ink object-cover';
        img.alt = '';
        img.src = URL.createObjectURL(file);
        img.onload = () => URL.revokeObjectURL(img.src);
        preview.append(img);
    });

    if (input.files.length > MAX_PHOTOS) {
        const note = document.createElement('p');
        note.className = 'col-span-full font-mono text-xs text-mute';
        note.textContent = `Only the first ${MAX_PHOTOS} photos will be sent.`;
        preview.append(note);
    }
}

/** After a successful create, move the new post into the feed and reset the form. */
function onPosted(response) {
    if (!response.querySelector('.post-success-marker')) return;
    const post = response.querySelector('.post');
    const feed = document.getElementById('posts-container');

    document.getElementById('create_post_modal')?.close();
    const form = document.getElementById('create_post_form');
    form?.reset();
    response.innerHTML = '';
    const preview = document.getElementById('image-preview');
    if (preview) {
        preview.innerHTML = '';
        preview.classList.add('hidden');
    }
    const editor = document.querySelector('#editor .ql-editor');
    if (editor) editor.innerHTML = '';

    if (feed && post) {
        feed.prepend(post);
        localizeDates(post);
        post.scrollIntoView({ behavior: 'smooth', block: 'start' });
    } else if (post) {
        window.location.href = `/doppler/#${post.id}`;
    }
}
