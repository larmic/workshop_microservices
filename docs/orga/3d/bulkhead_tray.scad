// Bulkhead tray for Story 5 (role play "Bulkhead").
// 3 columns (FLIGHT / HOTEL / CAR) x 2 wells for 40 mm poker chips.
// Round 1: no walls, one shared pool of 6 chips.
// Round 2: two plug-in walls between the columns, 2 chips per service.
//
// Print: tray upside up, walls lying flat. No supports needed.
// Export single parts with: openscad -D 'part="tray"' -o tray.stl bulkhead_tray.scad
// Tested with OpenSCAD 2021.01.

/* [Part] */
part = "all";            // "all" (print plate layout), "tray", "wall", "assembled" (preview only)

/* [Chips] */
chip_d    = 40;          // poker chip diameter
chip_h    = 3.3;         // poker chip thickness
stack     = 3;           // chips per well (stack height capacity)
clearance = 1.5;         // diameter clearance of a well

/* [Layout] */
columns   = 3;
rows      = 2;
col_pitch = 52;          // centre-to-centre distance between columns
row_gap   = 4.5;         // web between the wells of one column
rim_side  = 8;           // outer rim left/right
rim_back  = 6;           // outer rim at the back
rim_front = 16;          // front rim, carries the labels
floor_t   = 2;           // floor below the wells
corner_r  = 4;

/* [Finger notch] */
notch_d   = 22;          // scoop between the two wells of a column

/* [Walls] */
wall_t        = 3;       // wall thickness
wall_fit      = 0.3;     // total clearance wall <-> slot (thickness and length)
wall_above    = 25;      // visible wall height above the rim
slot_floor    = 2.5;     // material left below a slot
slot_end      = 3;       // closed slot ends (front/back)
slot_chamfer  = 0.6;     // lead-in chamfer at the slot opening
wall_text     = "BULKHEAD"; // debossed on the wall face, "" for none

/* [Labels] */
labels       = true;
label_texts  = ["FLIGHT", "HOTEL", "CAR"];
label_size   = 7;
label_depth  = 0.6;      // debossed into the front rim top
font         = "DejaVu Sans:style=Bold"; // e.g. "Liberation Sans:style=Bold"

$fn = 96;

// ---- derived ----
well_d     = chip_d + clearance;
well_h     = stack * chip_h + 0.5;
tray_h     = floor_t + well_h;
row_pitch  = well_d + row_gap;
tray_w     = (columns - 1) * col_pitch + well_d + 2 * rim_side;
tray_d     = rim_front + rim_back + (rows - 1) * row_pitch + well_d;
slot_w     = wall_t + wall_fit;
slot_len   = tray_d - 2 * slot_end;
slot_depth = tray_h - slot_floor;
wall_len   = slot_len - wall_fit;
wall_h     = slot_depth - 0.2 + wall_above;   // 0.2 mm air below the wall

function col_x(c) = rim_side + well_d / 2 + c * col_pitch;
function row_y(r) = rim_front + well_d / 2 + r * row_pitch;
function slot_x(i) = (col_x(i) + col_x(i + 1)) / 2;

echo(str("tray ", tray_w, " x ", tray_d, " x ", tray_h, " mm"));
echo(str("slot ", slot_w, " x ", slot_len, " x ", slot_depth, " mm"));
echo(str("wall ", wall_t, " x ", wall_len, " x ", wall_h, " mm"));
assert(col_pitch - well_d - slot_w >= 4, "too little material between well and slot");

module rounded_box(w, d, h, r) {
    translate([r, r, 0]) linear_extrude(h)
        offset(r = r) square([w - 2 * r, d - 2 * r]);
}

module tray() {
    difference() {
        rounded_box(tray_w, tray_d, tray_h, corner_r);

        // chip wells
        for (c = [0 : columns - 1], r = [0 : rows - 1])
            translate([col_x(c), row_y(r), floor_t]) cylinder(d = well_d, h = well_h + 1);

        // finger scoops between neighbouring wells of a column
        if (rows > 1)
            for (c = [0 : columns - 1], r = [0 : rows - 2])
                translate([col_x(c), (row_y(r) + row_y(r + 1)) / 2, floor_t])
                    cylinder(d = notch_d, h = well_h + 1);

        // wall slots with lead-in chamfer
        for (i = [0 : columns - 2])
            translate([slot_x(i), slot_end, 0]) {
                translate([-slot_w / 2, 0, slot_floor]) cube([slot_w, slot_len, slot_depth + 1]);
                translate([0, 0, tray_h - slot_chamfer]) rotate([-90, 0, 0])
                    linear_extrude(slot_len)
                        polygon([[-slot_w / 2 - slot_chamfer, -slot_chamfer - 0.01],
                                 [ slot_w / 2 + slot_chamfer, -slot_chamfer - 0.01],
                                 [ slot_w / 2, 0], [-slot_w / 2, 0]]);
            }

        // labels, debossed into the top of the front rim
        if (labels)
            for (c = [0 : min(columns, len(label_texts)) - 1])
                translate([col_x(c), rim_front / 2, tray_h - label_depth])
                    linear_extrude(label_depth + 1)
                        text(label_texts[c], size = label_size, font = font,
                             halign = "center", valign = "center");
    }
}

// Wall in print orientation: lying flat, x = length, y = height, z = thickness.
module wall() {
    c = 2; // chamfer of the upper corners
    difference() {
        linear_extrude(wall_t)
            polygon([[0, 0], [wall_len, 0], [wall_len, wall_h - c],
                     [wall_len - c, wall_h], [c, wall_h], [0, wall_h - c]]);
        if (wall_text != "")
            translate([wall_len / 2, slot_depth + wall_above / 2, wall_t - 0.6])
                linear_extrude(1)
                    text(wall_text, size = 8, font = font, halign = "center", valign = "center");
    }
}

if (part == "tray") tray();
else if (part == "wall") wall();
else if (part == "assembled") {
    tray();
    for (i = [0 : columns - 2])
        color("orange")
        translate([slot_x(i) + wall_t / 2, slot_end + wall_fit / 2, slot_floor + 0.2])
            rotate([90, 0, 90]) mirror([0, 0, 1]) wall();
} else {
    // print plate layout: tray plus all walls next to each other
    tray();
    for (i = [0 : columns - 2])
        translate([0, tray_d + 8 + i * (wall_h + 6), 0]) wall();
}
