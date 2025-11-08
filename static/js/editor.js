import Quill from 'quill';
import 'quill/dist/quill.snow.css';

export function initQuillEditor() {
    const modal = document.getElementById('create_post_modal');
    const form = document.querySelector('#create_post_form');
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

            console.log('Quill editor initialized');
        }
    };

    // Use HTMX's configRequest event to add Quill content before submission
    if (form) {
        form.addEventListener('htmx:configRequest', (event) => {
            if (quill) {
                const content = quill.root.innerHTML;
                // Add content to the HTMX request parameters
                event.detail.parameters['content'] = content;
                console.log('Quill content added to HTMX request:', content);
            }
        });
    }

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
}
