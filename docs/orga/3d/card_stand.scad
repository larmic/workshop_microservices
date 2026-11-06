// Stand for A6 landscape cards (148 x 105 mm, 250 g card stock).
// The slot leans back by `tilt` degrees. Print lying on the flat bottom,
// no supports needed. Set grid_x/grid_y > 1 to get a whole print plate.
// Tested with OpenSCAD 2021.01.

/* [Card] */
card_t      = 0.3;      // card stock thickness (250 g ~ 0.3 mm, laminated ~0.5-0.8 mm)
slot_extra  = 0.5;      // extra slot width so the card slides in easily
tilt        = 10;       // backward tilt in degrees

/* [Stand] */
len         = 70;       // along the card width
depth       = 40;       // front to back
front_h     = 6;        // height at the front edge
back_h      = 18;       // height at the slot
back_end_h  = 7;        // height at the back edge
slot_y      = 16;       // slot position from the front edge (at the top)
slot_depth  = 14;       // measured along the slot
corner_r    = 3;

/* [Plate] */
grid_x      = 1;
grid_y      = 1;
gap         = 6;

slot_w = card_t + slot_extra;
echo(str("stand ", len, " x ", depth, " x ", back_h, " mm, slot ", slot_w, " mm"));

module profile() {
    // side view in the y/z plane: low front, ramp up to the slot, ramp down behind it
    polygon([[0, 0], [depth, 0], [depth, back_end_h], [slot_y + 7, back_h],
             [slot_y - 4, back_h], [0, front_h]]);
}

module stand() {
    difference() {
        intersection() {
            rotate([90, 0, 90]) linear_extrude(len) profile();
            // round the corners in the footprint
            translate([corner_r, corner_r, -1]) linear_extrude(back_h + 2)
                offset(r = corner_r) square([len - 2 * corner_r, depth - 2 * corner_r]);
        }
        // slot, tilted backwards (top leans to +y)
        translate([-1, slot_y, back_h]) rotate([-tilt, 0, 0])
            translate([0, -slot_w / 2, -slot_depth]) cube([len + 2, slot_w, slot_depth + 10]);
        // small lead-in at the slot opening
        translate([-1, slot_y, back_h]) rotate([-tilt, 0, 0]) rotate([0, 90, 0])
            linear_extrude(len + 2) polygon([[-0.01, -slot_w / 2 - 0.8], [-0.01, slot_w / 2 + 0.8], [0.8, slot_w / 2], [0.8, -slot_w / 2]]);
    }
}

for (i = [0 : grid_x - 1], j = [0 : grid_y - 1])
    translate([i * (len + gap), j * (depth + gap), 0]) stand();
