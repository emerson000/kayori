'use client';

import { useState } from "react";
import { deleteProject } from "@/services/projectService";
import { useRouter } from "next/navigation";

interface ProjectData {
    id: string;
    title: string;
    status: string;
}

interface DeleteProjectButtonProps {
    project: ProjectData;
}

export default function DeleteProjectButton({ project }: DeleteProjectButtonProps) {
    const [showDeleteModal, setShowDeleteModal] = useState(false);
    const [isDeleting, setIsDeleting] = useState(false);
    const router = useRouter();

    const handleDelete = async () => {
        try {
            setIsDeleting(true);
            const success = await deleteProject(project.id);
            
            if (success) {
                // Redirect to projects list after successful deletion
                router.push('/projects');
            } else {
                // Handle failed deletion
                console.error('Failed to delete project');
            }
        } catch (error) {
            console.error('Error deleting project:', error);
        } finally {
            setIsDeleting(false);
            setShowDeleteModal(false);
        }
    };

    return (
        <>
            <button 
                className="btn btn-sm btn-outline btn-base-200 btn-error"
                onClick={() => setShowDeleteModal(true)}
            >
                Delete
            </button>

            {/* Delete Confirmation Modal */}
            <dialog className={`modal ${showDeleteModal ? 'modal-open' : ''}`}>
                <div className="modal-box">
                    <h3 className="font-bold text-lg">Delete Project</h3>
                    <p className="py-4">Are you sure you want to delete this project? This action cannot be undone.</p>
                    <div className="modal-action">
                        <button 
                            className="btn" 
                            onClick={() => setShowDeleteModal(false)}
                            disabled={isDeleting}
                        >
                            Cancel
                        </button>
                        <button 
                            className="btn btn-error" 
                            onClick={handleDelete}
                            disabled={isDeleting}
                        >
                            {isDeleting ? 'Deleting...' : 'Delete'}
                        </button>
                    </div>
                </div>
                <form method="dialog" className="modal-backdrop">
                    <button onClick={() => setShowDeleteModal(false)} disabled={isDeleting}>close</button>
                </form>
            </dialog>
        </>
    );
} 