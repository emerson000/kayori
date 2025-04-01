export const TLPLabel = {
    TLP_RED: { label: "TLP:RED", color: "bg-error text-error-content" },
    TLP_AMBER: { label: "TLP:AMBER", color: "bg-warning text-warning-content" },
    TLP_GREEN: { label: "TLP:GREEN", color: "bg-success text-success-content" },
    TLP_WHITE: { label: "TLP:WHITE", color: "bg-white text-black" },
} as const;

export type TLPLabelType = typeof TLPLabel[keyof typeof TLPLabel];

export default function Classification({ label }: { label: TLPLabelType }) {
    return <div>
        <div className={`${label.color} rounded-box text-center font-bold p-1 mb-4`}>
            {label.label}
        </div>
    </div>
}