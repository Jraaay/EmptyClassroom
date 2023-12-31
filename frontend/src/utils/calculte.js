export default CalculateEmptyClassroom;

function CalculateEmptyClassroom(classInfo, selectedCampus, selectedDate, selectedBuildings, selectedClassTimes) {
    let classInfoOnSelectedDate = classInfo.campus_info_map[selectedCampus];
    if (selectedDate.toDateString() != new Date().toDateString()) {
        classInfoOnSelectedDate = CalculateSelectedDateClassInfo(classInfo, selectedCampus, selectedDate);
    }
    let emptyClassroomList = [];
    for (let buildingId of selectedBuildings) {
        const buildingInfo = classInfoOnSelectedDate.building_info_map[buildingId];
        for (let classroomId in buildingInfo.classroom_info_map) {
            const classroomInfo = buildingInfo.classroom_info_map[classroomId];
            let emptyClassroom = {
                name: buildingInfo.name + '-' + classroomInfo.name,
                size: classroomInfo.size == 0 ? '无数据' : classroomInfo.size,
                can_trust: classroomInfo.can_trust,
                type: (classroomInfo.type == '' ? '无数据' : classroomInfo.type) ?? '无数据',
                empty_class_time: []
            };
            for (let classTime of selectedClassTimes) {
                if (buildingInfo.class_matrix[classTime][classroomId] == 1) {
                    emptyClassroom = null;
                    break;
                }
            }
            if (emptyClassroom != null) {
                for (let classTime = 0; classTime < 14; classTime++) {
                    if (buildingInfo.class_matrix[classTime][classroomId] == 0) {
                        emptyClassroom.empty_class_time.push(classTime);
                    }
                }
                emptyClassroomList.push(emptyClassroom);
            }
        }
    }
    emptyClassroomList.sort((a, b) => {
        if (a.building_ame != b.building_ame) {
            return a.building_ame > b.building_ame ? 1 : -1;
        } else {
            return a.name > b.name ? 1 : -1;
        }
    });
    return emptyClassroomList;
}

function CalculateSelectedDateClassInfo(classInfo, selectedCampus, selectedDate) {
    const buildingIdMap = classInfo.campus_info_map[selectedCampus].building_id_map;
    const classTable = classInfo.class_table.class_table_map[selectedCampus];
    const startWeek = Date.parse(classInfo.class_table.start_week);
    const nowWeek = Math.floor((selectedDate.getTime() - startWeek) / 604800000);
    const nowWeekDay = selectedDate.getDay();
    let resp = {
        building_id_map: buildingIdMap,
        building_info_map: {},
        max_building_id: 0
    };
    for (let classroomClass of classTable.class) {
        const buildingName = classroomClass.name.split('-')[0];
        const buildingId = buildingIdMap[buildingName];
        const classroomName = classroomClass.name.split('-')[1];
        if (resp.building_info_map[buildingId] == undefined) {
            resp.building_info_map[buildingId] = {
                name: buildingName,
                classroom_info_map: {},
                classroom_id_map: {},
                class_matrix: [],
                max_classroom_id: 0
            };
            for (let i = 0; i < 14; i++) {
                resp.building_info_map[buildingId].class_matrix[i] = [];
            }
        }
        const buildingInfo = resp.building_info_map[buildingId];
        if (buildingInfo.classroom_id_map[classroomName] == undefined) {
            buildingInfo.classroom_id_map[classroomName] = buildingInfo.max_classroom_id;
            buildingInfo.max_classroom_id += 1;
            buildingInfo.classroom_info_map[buildingInfo.classroom_id_map[classroomName]] = {
                name: classroomName,
                size: classroomClass.seat,
                can_trust: false,
                building_id: buildingId,
                type: classTable.typeMap[classroomClass.name]
            };
            for (let i = 0; i < 14; i++) {
                buildingInfo.class_matrix[i][buildingInfo.classroom_id_map[classroomName]] = 0;
            }
        }
        const classroomId = buildingInfo.classroom_id_map[classroomName];
        for (let i = 0; i < 14; i++) {
            if (classroomClass.classes[i][nowWeekDay] == null) {
                continue;
            }
            for (let week of classroomClass.classes[i][nowWeekDay]) {
                if (week == nowWeek) {
                    buildingInfo.class_matrix[i][classroomId] = 1;
                }
            }
        }
    }
    return resp;
}