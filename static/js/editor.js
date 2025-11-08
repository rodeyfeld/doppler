import Quill from 'quill';
import 'quill/dist/quill.snow.css';

export function initQuillEditor() {
    const modal = document.getElementById('create_post_modal');
    let quill = null;

    // Initialize Quill when modal is opened
    const initQuill = () => {
        if (!quill && document.getElementById('editor')) {
            quill = new Quill('#editor', {
                theme: 'snow',
                placeholder: 'Type your content here...',
                modules: {
                    toolbar: [
                        [{ 'header': [1, 2, 3, false] }],
                        ['bold', 'italic', 'underline', 'strike'],
                        ['blockquote', 'code-block'],
                        [{ 'list': 'ordered' }, { 'list': 'bullet' }],
                        ['link', 'image'],
                        ['clean']
                    ]
                }
            });
        }
    };

    // Initialize when modal opens
    if (modal) {
        modal.addEventListener('click', (e) => {
            if (modal.open) {
                initQuill();
            }
        });
    }

    // Initialize immediately if modal exists
    setTimeout(initQuill, 100);

    // Handle form submission to get Quill content
    const form = document.querySelector('#create_post_form');
    if (form) {
        form.addEventListener('submit', (e) => {
            if (quill) {
                const content = quill.root.innerHTML;
                document.querySelector('input[name="content"]').value = content;
            }
        });
    }
}

